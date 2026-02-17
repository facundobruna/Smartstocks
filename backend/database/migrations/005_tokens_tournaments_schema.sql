-- Smart Stocks Database Schema - MySQL
-- Fase 7: Sistema de Tokens y Torneos

-- ===========================================
-- TABLA: user_tokens (Moneda virtual)
-- ===========================================
CREATE TABLE user_tokens (
    user_id CHAR(36) PRIMARY KEY,
    balance INT DEFAULT 0,
    total_earned INT DEFAULT 0,
    total_spent INT DEFAULT 0,
    last_transaction_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tokens_balance (balance DESC),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: token_transactions (Historial de transacciones)
-- ===========================================
CREATE TABLE token_transactions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    transaction_type ENUM(
        'tournament_reward', 'tournament_entry',
        'daily_bonus', 'achievement_bonus',
        'admin_grant', 'purchase', 'refund'
    ) NOT NULL,
    amount INT NOT NULL,
    balance_after INT NOT NULL,
    description TEXT,
    reference_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_transactions_user (user_id),
    INDEX idx_transactions_type (transaction_type),
    INDEX idx_transactions_created (created_at DESC),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: tournaments (Torneos)
-- ===========================================
CREATE TABLE tournaments (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    tournament_type ENUM('weekly', 'monthly', 'special') NOT NULL,
    format ENUM('bracket', 'league', 'battle_royale') NOT NULL,
    entry_fee INT DEFAULT 0,
    prize_pool INT NOT NULL,
    min_rank_required VARCHAR(20) DEFAULT 'Bronze 1',
    max_participants INT NOT NULL,
    current_participants INT DEFAULT 0,
    status ENUM('upcoming', 'registration', 'in_progress', 'completed', 'cancelled') DEFAULT 'upcoming',
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    registration_start TIMESTAMP NOT NULL,
    registration_end TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tournaments_status (status),
    INDEX idx_tournaments_type (tournament_type),
    INDEX idx_tournaments_start (start_time),
    INDEX idx_tournaments_registration (registration_start, registration_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: tournament_participants (Participantes)
-- ===========================================
CREATE TABLE tournament_participants (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    tournament_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    current_score INT DEFAULT 0,
    current_position INT DEFAULT 0,
    matches_played INT DEFAULT 0,
    matches_won INT DEFAULT 0,
    matches_lost INT DEFAULT 0,
    is_eliminated BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_participants_tournament (tournament_id),
    INDEX idx_participants_user (user_id),
    INDEX idx_participants_score (tournament_id, current_score DESC),
    UNIQUE KEY unique_tournament_user (tournament_id, user_id),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: tournament_matches (Partidas del torneo)
-- ===========================================
CREATE TABLE tournament_matches (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    tournament_id CHAR(36) NOT NULL,
    round_number INT NOT NULL,
    match_number INT NOT NULL,
    player1_id CHAR(36) NOT NULL,
    player2_id CHAR(36) NOT NULL,
    player1_score INT DEFAULT 0,
    player2_score INT DEFAULT 0,
    winner_id CHAR(36),
    status ENUM('pending', 'in_progress', 'completed') DEFAULT 'pending',
    pvp_match_id CHAR(36),
    scheduled_time TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tournament_matches_tournament (tournament_id),
    INDEX idx_tournament_matches_round (tournament_id, round_number),
    INDEX idx_tournament_matches_status (status),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    FOREIGN KEY (player1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (player2_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (pvp_match_id) REFERENCES pvp_matches(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: tournament_prizes (Premios del torneo)
-- ===========================================
CREATE TABLE tournament_prizes (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    tournament_id CHAR(36) NOT NULL,
    position_from INT NOT NULL,
    position_to INT NOT NULL,
    token_reward INT NOT NULL,
    special_reward TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_prizes_tournament (tournament_id),
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TRIGGER: Crear balance de tokens al crear usuario
-- ===========================================
DELIMITER //
CREATE TRIGGER after_user_insert_tokens
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO user_tokens (user_id, balance, total_earned, total_spent)
    VALUES (NEW.id, 0, 0, 0);
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Añadir tokens
-- ===========================================
DELIMITER //
CREATE PROCEDURE add_tokens(
    IN p_user_id CHAR(36),
    IN p_amount INT,
    IN p_transaction_type VARCHAR(50),
    IN p_description TEXT,
    IN p_reference_id CHAR(36)
)
BEGIN
    DECLARE current_balance INT;
    DECLARE new_balance INT;

    -- Obtener balance actual
    SELECT balance INTO current_balance
    FROM user_tokens
    WHERE user_id = p_user_id;

    -- Calcular nuevo balance
    SET new_balance = current_balance + p_amount;

    -- Actualizar balance
    UPDATE user_tokens
    SET
        balance = new_balance,
        total_earned = total_earned + p_amount,
        last_transaction_at = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
    WHERE user_id = p_user_id;

    -- Registrar transacción
    INSERT INTO token_transactions (
        id, user_id, transaction_type, amount, balance_after,
        description, reference_id
    ) VALUES (
        UUID(), p_user_id, p_transaction_type, p_amount, new_balance,
        p_description, p_reference_id
    );
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Restar tokens
-- ===========================================
DELIMITER //
CREATE PROCEDURE subtract_tokens(
    IN p_user_id CHAR(36),
    IN p_amount INT,
    IN p_transaction_type VARCHAR(50),
    IN p_description TEXT,
    IN p_reference_id CHAR(36),
    OUT success BOOLEAN
)
BEGIN
    DECLARE current_balance INT;
    DECLARE new_balance INT;

    -- Obtener balance actual
    SELECT balance INTO current_balance
    FROM user_tokens
    WHERE user_id = p_user_id;

    -- Verificar si tiene suficientes tokens
    IF current_balance >= p_amount THEN
        SET new_balance = current_balance - p_amount;

        -- Actualizar balance
        UPDATE user_tokens
        SET
            balance = new_balance,
            total_spent = total_spent + p_amount,
            last_transaction_at = CURRENT_TIMESTAMP,
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = p_user_id;

        -- Registrar transacción (con cantidad negativa)
        INSERT INTO token_transactions (
            id, user_id, transaction_type, amount, balance_after,
            description, reference_id
        ) VALUES (
            UUID(), p_user_id, p_transaction_type, -p_amount, new_balance,
            p_description, p_reference_id
        );

        SET success = TRUE;
    ELSE
        SET success = FALSE;
    END IF;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Unirse a torneo
-- ===========================================
DELIMITER //
CREATE PROCEDURE join_tournament(
    IN p_tournament_id CHAR(36),
    IN p_user_id CHAR(36),
    OUT success BOOLEAN,
    OUT error_message VARCHAR(255)
)
BEGIN
    DECLARE entry_fee INT;
    DECLARE user_balance INT;
    DECLARE max_participants INT;
    DECLARE current_count INT;
    DECLARE tournament_status VARCHAR(20);
    DECLARE min_rank VARCHAR(20);
    DECLARE user_rank VARCHAR(20);
    DECLARE can_pay BOOLEAN;

    -- Obtener información del torneo
    SELECT t.entry_fee, t.max_participants, t.current_participants, t.status, t.min_rank_required
    INTO entry_fee, max_participants, current_count, tournament_status, min_rank
    FROM tournaments t
    WHERE t.id = p_tournament_id;

    -- Verificar que el torneo existe
    IF entry_fee IS NULL THEN
        SET success = FALSE;
        SET error_message = 'Tournament not found';
        LEAVE BEGIN;
    END IF;

    -- Verificar estado del torneo
    IF tournament_status != 'registration' THEN
        SET success = FALSE;
        SET error_message = 'Tournament is not in registration phase';
        LEAVE BEGIN;
    END IF;

    -- Verificar capacidad
    IF current_count >= max_participants THEN
        SET success = FALSE;
        SET error_message = 'Tournament is full';
        LEAVE BEGIN;
    END IF;

    -- Verificar rango mínimo
    SELECT us.rank_tier INTO user_rank
    FROM user_stats us
    WHERE us.user_id = p_user_id;

    IF NOT meets_rank_requirement(user_rank, min_rank) THEN
        SET success = FALSE;
        SET error_message = CONCAT('Minimum rank required: ', min_rank);
        LEAVE BEGIN;
    END IF;

    -- Verificar si ya está inscrito
    IF EXISTS (SELECT 1 FROM tournament_participants WHERE tournament_id = p_tournament_id AND user_id = p_user_id) THEN
        SET success = FALSE;
        SET error_message = 'Already registered in this tournament';
        LEAVE BEGIN;
    END IF;

    -- Cobrar entry fee si es necesario
    IF entry_fee > 0 THEN
        CALL subtract_tokens(
            p_user_id,
            entry_fee,
            'tournament_entry',
            CONCAT('Entry fee for tournament: ', p_tournament_id),
            p_tournament_id,
            can_pay
        );

        IF NOT can_pay THEN
            SET success = FALSE;
            SET error_message = 'Insufficient tokens';
            LEAVE BEGIN;
        END IF;
    END IF;

    -- Inscribir al usuario
    INSERT INTO tournament_participants (
        id, tournament_id, user_id, current_score, current_position
    ) VALUES (
        UUID(), p_tournament_id, p_user_id, 0, 0
    );

    -- Actualizar contador de participantes
    UPDATE tournaments
    SET current_participants = current_participants + 1
    WHERE id = p_tournament_id;

    SET success = TRUE;
    SET error_message = NULL;
END//
DELIMITER ;

-- ===========================================
-- FUNCIÓN: Verificar requisito de rango
-- ===========================================
DELIMITER //
CREATE FUNCTION meets_rank_requirement(
    user_rank VARCHAR(20),
    required_rank VARCHAR(20)
) RETURNS BOOLEAN
DETERMINISTIC
BEGIN
    DECLARE user_rank_value INT;
    DECLARE required_rank_value INT;

    -- Convertir rangos a valores numéricos
    SET user_rank_value = CASE
        WHEN user_rank LIKE 'Bronze%' THEN 1
        WHEN user_rank LIKE 'Plata%' THEN 2
        WHEN user_rank LIKE 'Oro%' THEN 3
        WHEN user_rank = 'Maestro' THEN 4
        ELSE 0
    END;

    SET required_rank_value = CASE
        WHEN required_rank LIKE 'Bronze%' THEN 1
        WHEN required_rank LIKE 'Plata%' THEN 2
        WHEN required_rank LIKE 'Oro%' THEN 3
        WHEN required_rank = 'Maestro' THEN 4
        ELSE 0
    END;

    RETURN user_rank_value >= required_rank_value;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Distribuir premios del torneo
-- ===========================================
DELIMITER //
CREATE PROCEDURE distribute_tournament_prizes(
    IN p_tournament_id CHAR(36)
)
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE v_user_id CHAR(36);
    DECLARE v_position INT;
    DECLARE v_token_reward INT;

    DECLARE prize_cursor CURSOR FOR
        SELECT
            tp.user_id,
            tp.current_position,
            pr.token_reward
        FROM tournament_participants tp
        JOIN tournament_prizes pr ON pr.tournament_id = tp.tournament_id
        WHERE tp.tournament_id = p_tournament_id
        AND tp.current_position >= pr.position_from
        AND tp.current_position <= pr.position_to
        ORDER BY tp.current_position ASC;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    OPEN prize_cursor;

    read_loop: LOOP
        FETCH prize_cursor INTO v_user_id, v_position, v_token_reward;
        IF done THEN
            LEAVE read_loop;
        END IF;

        -- Otorgar tokens
        CALL add_tokens(
            v_user_id,
            v_token_reward,
            'tournament_reward',
            CONCAT('Tournament prize - Position: ', v_position),
            p_tournament_id
        );
    END LOOP;

    CLOSE prize_cursor;

    -- Marcar torneo como completado
    UPDATE tournaments
    SET status = 'completed'
    WHERE id = p_tournament_id;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Actualizar posiciones del torneo
-- ===========================================
DELIMITER //
CREATE PROCEDURE update_tournament_positions(
    IN p_tournament_id CHAR(36)
)
BEGIN
    -- Actualizar posiciones basado en puntaje
    UPDATE tournament_participants tp
    JOIN (
        SELECT
            user_id,
            ROW_NUMBER() OVER (ORDER BY current_score DESC, joined_at ASC) as new_position
        FROM tournament_participants
        WHERE tournament_id = p_tournament_id
    ) ranked ON tp.user_id = ranked.user_id
    SET tp.current_position = ranked.new_position
    WHERE tp.tournament_id = p_tournament_id;
END//
DELIMITER ;

-- ===========================================
-- DATOS DE EJEMPLO: Torneo semanal
-- ===========================================
INSERT INTO tournaments (
    id, name, description, tournament_type, format,
    entry_fee, prize_pool, min_rank_required, max_participants,
    status, start_time, end_time, registration_start, registration_end
) VALUES (
    UUID(),
    'Torneo Semanal de Finanzas',
    'Compite contra los mejores traders cada semana',
    'weekly',
    'league',
    50,
    5000,
    'Plata 1',
    32,
    'registration',
    DATE_ADD(NOW(), INTERVAL 2 DAY),
    DATE_ADD(NOW(), INTERVAL 9 DAY),
    NOW(),
    DATE_ADD(NOW(), INTERVAL 1 DAY)
);

-- Agregar premios para el torneo de ejemplo
SET @tournament_id = (SELECT id FROM tournaments ORDER BY created_at DESC LIMIT 1);

INSERT INTO tournament_prizes (id, tournament_id, position_from, position_to, token_reward) VALUES
(UUID(), @tournament_id, 1, 1, 2000),
(UUID(), @tournament_id, 2, 2, 1000),
(UUID(), @tournament_id, 3, 3, 500),
(UUID(), @tournament_id, 4, 8, 250),
(UUID(), @tournament_id, 9, 16, 100);

-- ===========================================
-- ÍNDICES ADICIONALES PARA OPTIMIZACIÓN
-- ===========================================
CREATE INDEX idx_tournaments_active ON tournaments(status, start_time);
CREATE INDEX idx_participants_tournament_score ON tournament_participants(tournament_id, current_score DESC);
CREATE INDEX idx_token_transactions_user_date ON token_transactions(user_id, created_at DESC);