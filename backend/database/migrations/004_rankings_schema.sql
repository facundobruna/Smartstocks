-- Smart Stocks Database Schema - MySQL
-- Fase 6: Sistema de Rankings (Leaderboards)

-- ===========================================
-- TABLA: leaderboard_cache (Cache de rankings)
-- ===========================================
CREATE TABLE leaderboard_cache (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    cache_type ENUM('global', 'school') NOT NULL,
    school_id CHAR(36),
    rank_position INT NOT NULL,
    user_id CHAR(36) NOT NULL,
    username VARCHAR(50) NOT NULL,
    smartpoints INT NOT NULL,
    rank_tier VARCHAR(20) NOT NULL,
    total_wins INT DEFAULT 0,
    total_losses INT DEFAULT 0,
    win_rate DECIMAL(5,2) DEFAULT 0.00,
    profile_picture_url VARCHAR(500),
    school_name VARCHAR(255),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_leaderboard_type (cache_type),
    INDEX idx_leaderboard_school (school_id),
    INDEX idx_leaderboard_rank (rank_position),
    INDEX idx_leaderboard_user (user_id),
    UNIQUE KEY unique_cache_position (cache_type, school_id, rank_position),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: user_achievements (Logros de usuarios)
-- ===========================================
CREATE TABLE user_achievements (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    achievement_type ENUM(
        'first_win', 'win_streak_3', 'win_streak_5', 'win_streak_10',
        'rank_bronze', 'rank_silver', 'rank_gold', 'rank_master',
        'quiz_master', 'pvp_legend', 'simulator_expert',
        'points_1000', 'points_5000', 'points_10000',
        'perfect_quiz', 'speed_demon', 'comeback_king'
    ) NOT NULL,
    achievement_name VARCHAR(100) NOT NULL,
    achievement_description TEXT,
    icon_url VARCHAR(500),
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_achievements_user (user_id),
    INDEX idx_achievements_type (achievement_type),
    UNIQUE KEY unique_user_achievement (user_id, achievement_type),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- VISTA: Top 100 Global Leaderboard
-- ===========================================
CREATE OR REPLACE VIEW leaderboard_global_top100 AS
SELECT
    ROW_NUMBER() OVER (ORDER BY us.smartpoints DESC, u.created_at ASC) as rank_position,
    u.id as user_id,
    u.username,
    us.smartpoints,
    us.rank_tier,
    us.total_wins,
    us.total_losses,
    CASE
        WHEN (us.total_wins + us.total_losses) > 0
        THEN ROUND((us.total_wins * 100.0) / (us.total_wins + us.total_losses), 2)
        ELSE 0.00
    END as win_rate,
    u.profile_picture_url,
    s.name as school_name,
    s.id as school_id
FROM users u
JOIN user_stats us ON u.id = us.user_id
LEFT JOIN schools s ON u.school_id = s.id
ORDER BY us.smartpoints DESC, u.created_at ASC
LIMIT 100;

-- ===========================================
-- STORED PROCEDURE: Actualizar cache de rankings
-- ===========================================
DELIMITER //
CREATE PROCEDURE update_leaderboard_cache()
BEGIN
    -- Limpiar cache existente
    DELETE FROM leaderboard_cache;

    -- Insertar ranking global (top 1000)
    INSERT INTO leaderboard_cache (
        id, cache_type, school_id, rank_position, user_id, username,
        smartpoints, rank_tier, total_wins, total_losses, win_rate,
        profile_picture_url, school_name
    )
    SELECT
        UUID(),
        'global',
        NULL,
        ROW_NUMBER() OVER (ORDER BY us.smartpoints DESC, u.created_at ASC),
        u.id,
        u.username,
        us.smartpoints,
        us.rank_tier,
        us.total_wins,
        us.total_losses,
        CASE
            WHEN (us.total_wins + us.total_losses) > 0
            THEN ROUND((us.total_wins * 100.0) / (us.total_wins + us.total_losses), 2)
            ELSE 0.00
        END,
        u.profile_picture_url,
        s.name
    FROM users u
    JOIN user_stats us ON u.id = us.user_id
    LEFT JOIN schools s ON u.school_id = s.id
    ORDER BY us.smartpoints DESC, u.created_at ASC
    LIMIT 1000;

    -- Insertar rankings por colegio (top 100 por colegio)
    INSERT INTO leaderboard_cache (
        id, cache_type, school_id, rank_position, user_id, username,
        smartpoints, rank_tier, total_wins, total_losses, win_rate,
        profile_picture_url, school_name
    )
    SELECT
        UUID(),
        'school',
        u.school_id,
        ROW_NUMBER() OVER (PARTITION BY u.school_id ORDER BY us.smartpoints DESC, u.created_at ASC),
        u.id,
        u.username,
        us.smartpoints,
        us.rank_tier,
        us.total_wins,
        us.total_losses,
        CASE
            WHEN (us.total_wins + us.total_losses) > 0
            THEN ROUND((us.total_wins * 100.0) / (us.total_wins + us.total_losses), 2)
            ELSE 0.00
        END,
        u.profile_picture_url,
        s.name
    FROM users u
    JOIN user_stats us ON u.id = us.user_id
    JOIN schools s ON u.school_id = s.id
    WHERE u.school_id IS NOT NULL
    ORDER BY u.school_id, us.smartpoints DESC, u.created_at ASC;

END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Obtener posición del usuario
-- ===========================================
DELIMITER //
CREATE PROCEDURE get_user_position(
    IN p_user_id CHAR(36),
    OUT global_position INT,
    OUT school_position INT
)
BEGIN
    -- Posición global
    SELECT rank_position INTO global_position
    FROM (
        SELECT
            u.id,
            ROW_NUMBER() OVER (ORDER BY us.smartpoints DESC, u.created_at ASC) as rank_position
        FROM users u
        JOIN user_stats us ON u.id = us.user_id
    ) ranks
    WHERE id = p_user_id;

    -- Posición en colegio
    SELECT rank_position INTO school_position
    FROM (
        SELECT
            u.id,
            ROW_NUMBER() OVER (PARTITION BY u.school_id ORDER BY us.smartpoints DESC, u.created_at ASC) as rank_position
        FROM users u
        JOIN user_stats us ON u.id = us.user_id
        WHERE u.school_id IS NOT NULL
    ) school_ranks
    WHERE id = p_user_id;

    -- Si no tiene colegio, posición es NULL
    IF school_position IS NULL THEN
        SET school_position = 0;
    END IF;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Otorgar logro
-- ===========================================
DELIMITER //
CREATE PROCEDURE grant_achievement(
    IN p_user_id CHAR(36),
    IN p_achievement_type VARCHAR(50),
    IN p_achievement_name VARCHAR(100),
    IN p_achievement_description TEXT
)
BEGIN
    -- Insertar solo si no existe
    INSERT IGNORE INTO user_achievements (
        id, user_id, achievement_type, achievement_name, achievement_description
    ) VALUES (
        UUID(), p_user_id, p_achievement_type, p_achievement_name, p_achievement_description
    );
END//
DELIMITER ;

-- ===========================================
-- TRIGGER: Otorgar logros automáticamente
-- ===========================================
DELIMITER //
CREATE TRIGGER check_achievements_after_stats_update
AFTER UPDATE ON user_stats
FOR EACH ROW
BEGIN
    -- Logro: Primera victoria
    IF NEW.total_wins = 1 AND OLD.total_wins = 0 THEN
        CALL grant_achievement(
            NEW.user_id,
            'first_win',
            'Primera Victoria',
            'Ganaste tu primera partida PvP'
        );
    END IF;

    -- Logro: Racha de 3 victorias
    IF NEW.win_streak = 3 AND OLD.win_streak < 3 THEN
        CALL grant_achievement(
            NEW.user_id,
            'win_streak_3',
            'En Racha',
            '3 victorias seguidas'
        );
    END IF;

    -- Logro: Racha de 5 victorias
    IF NEW.win_streak = 5 AND OLD.win_streak < 5 THEN
        CALL grant_achievement(
            NEW.user_id,
            'win_streak_5',
            'Imparable',
            '5 victorias seguidas'
        );
    END IF;

    -- Logro: Racha de 10 victorias
    IF NEW.win_streak = 10 AND OLD.win_streak < 10 THEN
        CALL grant_achievement(
            NEW.user_id,
            'win_streak_10',
            'Leyenda',
            '10 victorias seguidas'
        );
    END IF;

    -- Logro: Alcanzar Plata
    IF NEW.rank_tier LIKE 'Plata%' AND OLD.rank_tier LIKE 'Bronze%' THEN
        CALL grant_achievement(
            NEW.user_id,
            'rank_silver',
            'Ascenso a Plata',
            'Alcanzaste el rango Plata'
        );
    END IF;

    -- Logro: Alcanzar Oro
    IF NEW.rank_tier LIKE 'Oro%' AND OLD.rank_tier LIKE 'Plata%' THEN
        CALL grant_achievement(
            NEW.user_id,
            'rank_gold',
            'Ascenso a Oro',
            'Alcanzaste el rango Oro'
        );
    END IF;

    -- Logro: Alcanzar Maestro
    IF NEW.rank_tier = 'Maestro' AND OLD.rank_tier LIKE 'Oro%' THEN
        CALL grant_achievement(
            NEW.user_id,
            'rank_master',
            'Maestro de las Finanzas',
            'Alcanzaste el rango Maestro'
        );
    END IF;

    -- Logro: 1,000 puntos
    IF NEW.smartpoints >= 1000 AND OLD.smartpoints < 1000 THEN
        CALL grant_achievement(
            NEW.user_id,
            'points_1000',
            'Mil Puntos',
            'Alcanzaste 1,000 SmartPoints'
        );
    END IF;

    -- Logro: 5,000 puntos
    IF NEW.smartpoints >= 5000 AND OLD.smartpoints < 5000 THEN
        CALL grant_achievement(
            NEW.user_id,
            'points_5000',
            'Cinco Mil',
            'Alcanzaste 5,000 SmartPoints'
        );
    END IF;

    -- Logro: 10,000 puntos
    IF NEW.smartpoints >= 10000 AND OLD.smartpoints < 10000 THEN
        CALL grant_achievement(
            NEW.user_id,
            'points_10000',
            'Diez Mil',
            'Alcanzaste 10,000 SmartPoints'
        );
    END IF;

    -- Logro: Quiz Master (50 quizzes completados)
    IF NEW.total_quizzes_completed >= 50 AND OLD.total_quizzes_completed < 50 THEN
        CALL grant_achievement(
            NEW.user_id,
            'quiz_master',
            'Maestro de Quizzes',
            'Completaste 50 quizzes'
        );
    END IF;
END//
DELIMITER ;

-- ===========================================
-- EVENT: Actualizar cache cada 5 minutos
-- ===========================================
SET GLOBAL event_scheduler = ON;

CREATE EVENT IF NOT EXISTS update_leaderboard_cache_event
ON SCHEDULE EVERY 5 MINUTE
DO
CALL update_leaderboard_cache();

-- ===========================================
-- Ejecutar actualización inicial
-- ===========================================
CALL update_leaderboard_cache();

-- ===========================================
-- ÍNDICES ADICIONALES PARA OPTIMIZACIÓN
-- ===========================================
CREATE INDEX idx_user_stats_points_desc ON user_stats(smartpoints DESC);
CREATE INDEX idx_users_school ON users(school_id);
CREATE INDEX idx_leaderboard_cache_updated ON leaderboard_cache(last_updated DESC);