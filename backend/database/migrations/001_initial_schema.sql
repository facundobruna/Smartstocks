-- Smart Stocks Database Schema - MySQL
-- Fase 1: Usuarios, Autenticación, Colegios, Stats Básicas

-- ===========================================
-- TABLA: schools (Colegios asociados)
-- ===========================================
CREATE TABLE schools (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_schools_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: users (Usuarios principales)
-- ===========================================
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    profile_picture_url VARCHAR(500),
    school_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    reset_token VARCHAR(255),
    reset_token_expires TIMESTAMP NULL,
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_school (school_id),
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: user_stats (Estadísticas de usuario)
-- ===========================================
CREATE TABLE user_stats (
    user_id CHAR(36) PRIMARY KEY,
    smartpoints INT DEFAULT 0,
    rank_tier VARCHAR(20) DEFAULT 'Bronze 1',
    total_quizzes_completed INT DEFAULT 0,
    total_simulator_games INT DEFAULT 0,
    win_streak INT DEFAULT 0,
    total_wins INT DEFAULT 0,
    total_losses INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_stats_points (smartpoints DESC),
    INDEX idx_user_stats_rank (rank_tier),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: refresh_tokens (Tokens de sesión)
-- ===========================================
CREATE TABLE refresh_tokens (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    token VARCHAR(500) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_refresh_tokens_user (user_id),
    INDEX idx_refresh_tokens_token (token),
    INDEX idx_refresh_tokens_expires (expires_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- INSERTAR COLEGIOS DE EJEMPLO
-- ===========================================
INSERT INTO schools (id, name, location, is_active) VALUES
(UUID(), 'Colegio Nacional Buenos Aires', 'Buenos Aires, Argentina', TRUE),
(UUID(), 'Instituto San Martín', 'Córdoba, Argentina', TRUE),
(UUID(), 'Escuela Técnica N°1', 'Rosario, Argentina', TRUE),
(UUID(), 'Colegio Belgrano', 'Mendoza, Argentina', TRUE),
(UUID(), 'Instituto Comercial', 'La Plata, Argentina', TRUE);

-- ===========================================
-- TRIGGER: Crear user_stats automáticamente
-- ===========================================
DELIMITER //
CREATE TRIGGER after_user_insert
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO user_stats (user_id, smartpoints, rank_tier)
    VALUES (NEW.id, 0, 'Bronze 1');
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Actualizar rango de usuario
-- ===========================================
DELIMITER //
CREATE PROCEDURE update_user_rank(IN p_user_id CHAR(36))
BEGIN
    DECLARE v_points INT;
    DECLARE v_new_rank VARCHAR(20);

    SELECT smartpoints INTO v_points
    FROM user_stats
    WHERE user_id = p_user_id;

    SET v_new_rank = CASE
        WHEN v_points < 400 THEN 'Bronze 1'
        WHEN v_points < 1600 THEN 'Bronze 2'
        WHEN v_points < 3200 THEN 'Bronze 3'
        WHEN v_points < 4400 THEN 'Plata 1'
        WHEN v_points < 6400 THEN 'Plata 2'
        WHEN v_points < 8400 THEN 'Plata 3'
        WHEN v_points < 10000 THEN 'Oro 1'
        WHEN v_points < 12400 THEN 'Oro 2'
        WHEN v_points < 14400 THEN 'Oro 3'
        ELSE 'Maestro'
    END;

    UPDATE user_stats
    SET rank_tier = v_new_rank
    WHERE user_id = p_user_id;
END//
DELIMITER ;