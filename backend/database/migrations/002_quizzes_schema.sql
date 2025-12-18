-- Fase 2: Sistema de Quizzes
-- Migraciones para tablas de quizzes, preguntas y attemps

-- ===========================================
-- TABLA: quizzes (Quizzes disponibles)
-- ===========================================
CREATE TABLE quizzes (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    points_reward INT NOT NULL,
    total_questions INT DEFAULT 10,
    time_limit_minutes INT DEFAULT 30,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    INDEX idx_quizzes_difficulty (difficulty),
    INDEX idx_quizzes_active (is_active),
    INDEX idx_quizzes_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: quiz_questions (Preguntas de quizzes)
-- ===========================================
CREATE TABLE quiz_questions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    quiz_id CHAR(36) NOT NULL,
    question_text TEXT NOT NULL,
    option_a TEXT NOT NULL,
    option_b TEXT NOT NULL,
    option_c TEXT NOT NULL,
    option_d TEXT NOT NULL,
    correct_option ENUM('A', 'B', 'C', 'D') NOT NULL,
    explanation TEXT,
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_questions_quiz (quiz_id),
    INDEX idx_questions_difficulty (difficulty),
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: quiz_attempts (Intentos de usuarios)
-- ===========================================
CREATE TABLE quiz_attempts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    quiz_id CHAR(36) NOT NULL,
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    score INT NOT NULL,
    total_questions INT NOT NULL,
    correct_answers INT NOT NULL,
    points_earned INT NOT NULL,
    time_taken_seconds INT,
    answers JSON,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_attempts_user (user_id),
    INDEX idx_attempts_quiz (quiz_id),
    INDEX idx_attempts_date (completed_at),
    UNIQUE KEY unique_user_quiz_date (user_id, quiz_id, DATE(completed_at)),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: daily_quiz_cooldowns (Control de cooldown)
-- ===========================================
CREATE TABLE daily_quiz_cooldowns (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    last_attempt_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_difficulty_date (user_id, difficulty, last_attempt_date),
    INDEX idx_cooldown_user (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- STORED PROCEDURE: Actualizar stats después de quiz
-- ===========================================
DELIMITER //
CREATE PROCEDURE update_user_stats_after_quiz(
    IN p_user_id CHAR(36),
    IN p_points_earned INT
)
BEGIN
    -- Actualizar smartpoints y total de quizzes
    UPDATE user_stats
    SET smartpoints = smartpoints + p_points_earned,
        total_quizzes_completed = total_quizzes_completed + 1,
        updated_at = NOW()
    WHERE user_id = p_user_id;

    -- Actualizar rango
    CALL update_user_rank(p_user_id);
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Verificar si puede hacer quiz
-- ===========================================
DELIMITER //
CREATE PROCEDURE check_quiz_cooldown(
    IN p_user_id CHAR(36),
    IN p_difficulty VARCHAR(10),
    OUT p_can_attempt BOOLEAN
)
BEGIN
    DECLARE v_count INT;

    SELECT COUNT(*) INTO v_count
    FROM daily_quiz_cooldowns
    WHERE user_id = p_user_id
    AND difficulty = p_difficulty
    AND last_attempt_date = CURDATE();

    SET p_can_attempt = (v_count = 0);
END//
DELIMITER ;

-- ===========================================
-- FUNCIÓN: Obtener puntos por dificultad
-- ===========================================
DELIMITER //
CREATE FUNCTION get_quiz_points(p_difficulty VARCHAR(10))
RETURNS INT
DETERMINISTIC
BEGIN
    RETURN CASE p_difficulty
        WHEN 'easy' THEN 500
        WHEN 'medium' THEN 1000
        WHEN 'hard' THEN 2000
        ELSE 0
    END;
END//
DELIMITER ;