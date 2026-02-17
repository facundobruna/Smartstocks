-- Fase 4: Sistema de Simulador de Trading
-- Migraciones para escenarios, intentos y cooldowns

-- ===========================================
-- TABLA: simulator_scenarios (Escenarios de trading)
-- ===========================================
CREATE TABLE IF NOT EXISTS simulator_scenarios (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    news_content TEXT NOT NULL,
    chart_data JSON NOT NULL,
    correct_decision ENUM('buy', 'sell', 'hold') NOT NULL,
    explanation TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    INDEX idx_scenarios_difficulty (difficulty),
    INDEX idx_scenarios_active (is_active),
    INDEX idx_scenarios_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: simulator_attempts (Intentos de usuarios)
-- ===========================================
CREATE TABLE IF NOT EXISTS simulator_attempts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    scenario_id CHAR(36) NOT NULL,
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    user_decision ENUM('buy', 'sell', 'hold') NOT NULL,
    was_correct BOOLEAN NOT NULL,
    points_earned INT DEFAULT 0,
    time_taken_seconds INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_attempts_user (user_id),
    INDEX idx_attempts_scenario (scenario_id),
    INDEX idx_attempts_difficulty (difficulty),
    INDEX idx_attempts_date (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (scenario_id) REFERENCES simulator_scenarios(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: daily_simulator_cooldowns (Control de cooldown)
-- ===========================================
CREATE TABLE IF NOT EXISTS daily_simulator_cooldowns (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    last_attempt_date DATE NOT NULL,
    attempts_count INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_difficulty_date (user_id, difficulty, last_attempt_date),
    INDEX idx_cooldown_user (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- STORED PROCEDURE: Verificar cooldown del simulador
-- ===========================================
DELIMITER //
CREATE PROCEDURE check_simulator_cooldown(
    IN p_user_id CHAR(36),
    IN p_difficulty VARCHAR(10),
    OUT p_can_attempt BOOLEAN
)
BEGIN
    DECLARE v_count INT;

    SELECT COUNT(*) INTO v_count
    FROM daily_simulator_cooldowns
    WHERE user_id = p_user_id
    AND difficulty = p_difficulty
    AND last_attempt_date = CURDATE();

    SET p_can_attempt = (v_count = 0);
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Registrar intento de simulador
-- ===========================================
DELIMITER //
CREATE PROCEDURE record_simulator_attempt(
    IN p_user_id CHAR(36),
    IN p_scenario_id CHAR(36),
    IN p_difficulty VARCHAR(10),
    IN p_user_decision VARCHAR(10),
    IN p_was_correct BOOLEAN,
    IN p_points_earned INT,
    IN p_time_taken INT
)
BEGIN
    DECLARE v_attempt_id CHAR(36);
    SET v_attempt_id = UUID();

    -- Insertar intento
    INSERT INTO simulator_attempts (
        id, user_id, scenario_id, difficulty,
        user_decision, was_correct, points_earned, time_taken_seconds
    ) VALUES (
        v_attempt_id, p_user_id, p_scenario_id, p_difficulty,
        p_user_decision, p_was_correct, p_points_earned, p_time_taken
    );

    -- Registrar cooldown
    INSERT INTO daily_simulator_cooldowns (id, user_id, difficulty, last_attempt_date)
    VALUES (UUID(), p_user_id, p_difficulty, CURDATE())
    ON DUPLICATE KEY UPDATE
        attempts_count = attempts_count + 1,
        updated_at = NOW();

    -- Si fue correcto, actualizar stats del usuario
    IF p_was_correct THEN
        UPDATE user_stats
        SET smartpoints = smartpoints + p_points_earned,
            total_simulator_games = total_simulator_games + 1,
            updated_at = NOW()
        WHERE user_id = p_user_id;

        -- Actualizar rango
        CALL update_user_rank(p_user_id);
    ELSE
        -- Solo incrementar contador de simulaciones
        UPDATE user_stats
        SET total_simulator_games = total_simulator_games + 1,
            updated_at = NOW()
        WHERE user_id = p_user_id;
    END IF;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Limpiar escenarios expirados
-- ===========================================
DELIMITER //
CREATE PROCEDURE cleanup_expired_scenarios()
BEGIN
    UPDATE simulator_scenarios
    SET is_active = FALSE
    WHERE expires_at < NOW() AND is_active = TRUE;
END//
DELIMITER ;

-- ===========================================
-- FUNCION: Obtener puntos por dificultad del simulador
-- ===========================================
DELIMITER //
CREATE FUNCTION get_simulator_points(p_difficulty VARCHAR(10))
RETURNS INT
DETERMINISTIC
BEGIN
    RETURN CASE p_difficulty
        WHEN 'easy' THEN 25
        WHEN 'medium' THEN 50
        WHEN 'hard' THEN 100
        ELSE 0
    END;
END//
DELIMITER ;
