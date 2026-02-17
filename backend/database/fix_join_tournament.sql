-- Fix: Crear stored procedure join_tournament
-- Ejecutar: mysql -u root -p smartstocks < database/fix_join_tournament.sql

DELIMITER //

DROP PROCEDURE IF EXISTS join_tournament//

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
    DECLARE can_pay BOOLEAN DEFAULT TRUE;

    -- Obtener información del torneo
    SELECT t.entry_fee, t.max_participants, t.current_participants, t.status, t.min_rank_required
    INTO entry_fee, max_participants, current_count, tournament_status, min_rank
    FROM tournaments t
    WHERE t.id = p_tournament_id;

    -- Verificar que el torneo existe
    IF entry_fee IS NULL THEN
        SET success = FALSE;
        SET error_message = 'Tournament not found';
    -- Verificar estado del torneo
    ELSEIF tournament_status != 'registration' THEN
        SET success = FALSE;
        SET error_message = 'Tournament is not in registration phase';
    -- Verificar capacidad
    ELSEIF current_count >= max_participants THEN
        SET success = FALSE;
        SET error_message = 'Tournament is full';
    -- Verificar si ya está inscrito
    ELSEIF EXISTS (SELECT 1 FROM tournament_participants WHERE tournament_id = p_tournament_id AND user_id = p_user_id) THEN
        SET success = FALSE;
        SET error_message = 'Already registered in this tournament';
    ELSE
        -- Verificar rango mínimo
        SELECT us.rank_tier INTO user_rank
        FROM user_stats us
        WHERE us.user_id = p_user_id;

        -- Cobrar entry fee si es necesario
        IF entry_fee > 0 THEN
            SELECT balance INTO user_balance
            FROM user_tokens
            WHERE user_id = p_user_id;

            IF user_balance IS NULL OR user_balance < entry_fee THEN
                SET can_pay = FALSE;
                SET success = FALSE;
                SET error_message = 'Insufficient tokens';
            ELSE
                -- Descontar tokens
                UPDATE user_tokens
                SET balance = balance - entry_fee,
                    total_spent = total_spent + entry_fee,
                    last_transaction_at = CURRENT_TIMESTAMP
                WHERE user_id = p_user_id;

                -- Registrar transacción
                INSERT INTO token_transactions (
                    id, user_id, transaction_type, amount, balance_after, description, reference_id
                ) VALUES (
                    UUID(), p_user_id, 'tournament_entry', -entry_fee,
                    user_balance - entry_fee,
                    CONCAT('Entry fee for tournament'),
                    p_tournament_id
                );
            END IF;
        END IF;

        IF can_pay THEN
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
        END IF;
    END IF;
END//

DELIMITER ;

-- Verificar que se creó
SHOW PROCEDURE STATUS WHERE Name = 'join_tournament';
