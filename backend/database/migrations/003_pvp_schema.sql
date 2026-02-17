-- Smart Stocks Database Schema - MySQL
-- Fase 3: Sistema PvP

-- ===========================================
-- TABLA: pvp_queue (Cola de matchmaking)
-- ===========================================
CREATE TABLE pvp_queue (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL UNIQUE,
    rank_tier VARCHAR(20) NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    INDEX idx_queue_active (is_active, joined_at),
    INDEX idx_queue_rank (rank_tier, is_active),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: pvp_matches (Partidas PvP)
-- ===========================================
CREATE TABLE pvp_matches (
    id CHAR(36) PRIMARY KEY,
    player1_id CHAR(36) NOT NULL,
    player2_id CHAR(36) NOT NULL,
    player1_score INT DEFAULT 0,
    player2_score INT DEFAULT 0,
    winner_id CHAR(36),
    status ENUM('waiting', 'in_progress', 'completed', 'abandoned') DEFAULT 'waiting',
    current_round INT DEFAULT 0,
    total_rounds INT DEFAULT 5,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_matches_player1 (player1_id),
    INDEX idx_matches_player2 (player2_id),
    INDEX idx_matches_status (status),
    INDEX idx_matches_created (created_at DESC),
    FOREIGN KEY (player1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (player2_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: pvp_rounds (Rondas de partida PvP)
-- ===========================================
CREATE TABLE pvp_rounds (
    id CHAR(36) PRIMARY KEY,
    match_id CHAR(36) NOT NULL,
    round_number INT NOT NULL,
    scenario_id CHAR(36) NOT NULL,
    player1_decision VARCHAR(10),
    player2_decision VARCHAR(10),
    player1_time_seconds DECIMAL(5,2),
    player2_time_seconds DECIMAL(5,2),
    player1_correct BOOLEAN,
    player2_correct BOOLEAN,
    player1_points INT DEFAULT 0,
    player2_points INT DEFAULT 0,
    correct_decision VARCHAR(10) NOT NULL,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_rounds_match (match_id),
    INDEX idx_rounds_match_number (match_id, round_number),
    UNIQUE KEY unique_match_round (match_id, round_number),
    FOREIGN KEY (match_id) REFERENCES pvp_matches(id) ON DELETE CASCADE,
    FOREIGN KEY (scenario_id) REFERENCES simulator_scenarios(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- STORED PROCEDURE: Buscar oponente
-- ===========================================
DELIMITER //
CREATE PROCEDURE find_opponent(
    IN p_user_id CHAR(36),
    IN p_rank_tier VARCHAR(20),
    OUT opponent_id CHAR(36)
)
BEGIN
    DECLARE found_opponent CHAR(36);

    -- Buscar oponente en la misma tier primero
    SELECT user_id INTO found_opponent
    FROM pvp_queue
    WHERE user_id != p_user_id
      AND is_active = TRUE
      AND expires_at > NOW()
      AND rank_tier = p_rank_tier
    ORDER BY joined_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED;

    -- Si no hay oponente en la misma tier, buscar en cualquier tier
    IF found_opponent IS NULL THEN
        SELECT user_id INTO found_opponent
        FROM pvp_queue
        WHERE user_id != p_user_id
          AND is_active = TRUE
          AND expires_at > NOW()
        ORDER BY joined_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED;
    END IF;

    -- Si encontramos oponente, marcarlo como inactivo
    IF found_opponent IS NOT NULL THEN
        UPDATE pvp_queue
        SET is_active = FALSE
        WHERE user_id = found_opponent;
    END IF;

    SET opponent_id = found_opponent;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Limpiar cola expirada
-- ===========================================
DELIMITER //
CREATE PROCEDURE cleanup_expired_queue()
BEGIN
    DELETE FROM pvp_queue
    WHERE expires_at < NOW() OR is_active = FALSE;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Actualizar stats PvP
-- ===========================================
DELIMITER //
CREATE PROCEDURE update_pvp_stats(
    IN p_winner_id CHAR(36),
    IN p_loser_id CHAR(36),
    IN p_winner_points INT,
    IN p_is_win BOOLEAN
)
BEGIN
    -- Actualizar stats del ganador
    UPDATE user_stats
    SET
        smartpoints = smartpoints + p_winner_points,
        total_wins = total_wins + 1,
        win_streak = win_streak + 1,
        rank_tier = calculate_rank_tier(smartpoints + p_winner_points),
        updated_at = CURRENT_TIMESTAMP
    WHERE user_id = p_winner_id;

    -- Actualizar stats del perdedor
    UPDATE user_stats
    SET
        smartpoints = GREATEST(0, smartpoints - 100),
        total_losses = total_losses + 1,
        win_streak = 0,
        rank_tier = calculate_rank_tier(GREATEST(0, smartpoints - 100)),
        updated_at = CURRENT_TIMESTAMP
    WHERE user_id = p_loser_id;
END//
DELIMITER ;

-- ===========================================
-- FUNCION: Calcular tier de rango
-- ===========================================
DELIMITER //
CREATE FUNCTION calculate_rank_tier(points INT)
RETURNS VARCHAR(20)
DETERMINISTIC
BEGIN
    RETURN CASE
        WHEN points >= 10000 THEN 'Maestro'
        WHEN points >= 7500 THEN 'Oro 1'
        WHEN points >= 5000 THEN 'Oro 2'
        WHEN points >= 3500 THEN 'Oro 3'
        WHEN points >= 2500 THEN 'Plata 1'
        WHEN points >= 1500 THEN 'Plata 2'
        WHEN points >= 1000 THEN 'Plata 3'
        WHEN points >= 500 THEN 'Bronce 1'
        WHEN points >= 250 THEN 'Bronce 2'
        ELSE 'Bronce 3'
    END;
END//
DELIMITER ;

-- ===========================================
-- TRIGGER: Limpiar cola al completar match
-- ===========================================
DELIMITER //
CREATE TRIGGER after_match_complete
AFTER UPDATE ON pvp_matches
FOR EACH ROW
BEGIN
    IF NEW.status = 'completed' AND OLD.status != 'completed' THEN
        -- Limpiar entradas de cola de ambos jugadores
        DELETE FROM pvp_queue WHERE user_id IN (NEW.player1_id, NEW.player2_id);
    END IF;
END//
DELIMITER ;
