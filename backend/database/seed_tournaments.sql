-- Seed de Torneos para SmartStocks
-- Ejecutar: mysql -u root -p smartstocks < database/seed_tournaments.sql

-- Limpiar torneos existentes (opcional)
-- DELETE FROM tournament_prizes;
-- DELETE FROM tournament_participants;
-- DELETE FROM tournaments;

-- ===========================================
-- TORNEO 1: Torneo Semanal (Inscripción Abierta)
-- ===========================================
SET @torneo1_id = UUID();

INSERT INTO tournaments (
    id, name, description, tournament_type, format,
    entry_fee, prize_pool, min_rank_required, max_participants,
    status, start_time, end_time, registration_start, registration_end
) VALUES (
    @torneo1_id,
    'Torneo Semanal - Traders Novatos',
    'Compite contra otros traders principiantes. Perfecto para empezar en el mundo de los torneos.',
    'weekly',
    'league',
    0,
    1000,
    'Bronce 3',
    16,
    'registration',
    DATE_ADD(NOW(), INTERVAL 2 DAY),
    DATE_ADD(NOW(), INTERVAL 5 DAY),
    DATE_SUB(NOW(), INTERVAL 1 DAY),
    DATE_ADD(NOW(), INTERVAL 1 DAY)
);

INSERT INTO tournament_prizes (id, tournament_id, position_from, position_to, token_reward, special_reward) VALUES
(UUID(), @torneo1_id, 1, 1, 500, 'Badge Campeón Semanal'),
(UUID(), @torneo1_id, 2, 2, 300, NULL),
(UUID(), @torneo1_id, 3, 3, 200, NULL);

-- ===========================================
-- TORNEO 2: Torneo Mensual Pro (Inscripción Abierta)
-- ===========================================
SET @torneo2_id = UUID();

INSERT INTO tournaments (
    id, name, description, tournament_type, format,
    entry_fee, prize_pool, min_rank_required, max_participants,
    status, start_time, end_time, registration_start, registration_end
) VALUES (
    @torneo2_id,
    'Copa Mensual de Inversores',
    'El torneo más prestigioso del mes. Solo para traders experimentados con rango Plata o superior.',
    'monthly',
    'bracket',
    100,
    10000,
    'Plata 3',
    32,
    'registration',
    DATE_ADD(NOW(), INTERVAL 5 DAY),
    DATE_ADD(NOW(), INTERVAL 12 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_ADD(NOW(), INTERVAL 3 DAY)
);

INSERT INTO tournament_prizes (id, tournament_id, position_from, position_to, token_reward, special_reward) VALUES
(UUID(), @torneo2_id, 1, 1, 5000, 'Título: Gran Maestro del Mes'),
(UUID(), @torneo2_id, 2, 2, 2500, 'Badge Subcampeón'),
(UUID(), @torneo2_id, 3, 3, 1500, 'Badge Tercer Lugar'),
(UUID(), @torneo2_id, 4, 8, 500, NULL),
(UUID(), @torneo2_id, 9, 16, 200, NULL);

-- ===========================================
-- TORNEO 3: Torneo Especial (En Progreso)
-- ===========================================
SET @torneo3_id = UUID();

INSERT INTO tournaments (
    id, name, description, tournament_type, format,
    entry_fee, prize_pool, min_rank_required, max_participants,
    current_participants, status, start_time, end_time, registration_start, registration_end
) VALUES (
    @torneo3_id,
    'Battle Royale: Supervivencia Financiera',
    'Un torneo especial donde solo el más astuto sobrevive. 50 participantes, 1 ganador.',
    'special',
    'battle_royale',
    50,
    7500,
    'Bronce 1',
    50,
    48,
    'in_progress',
    DATE_SUB(NOW(), INTERVAL 1 DAY),
    DATE_ADD(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 5 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY)
);

INSERT INTO tournament_prizes (id, tournament_id, position_from, position_to, token_reward, special_reward) VALUES
(UUID(), @torneo3_id, 1, 1, 4000, 'Título: Último en Pie'),
(UUID(), @torneo3_id, 2, 2, 2000, NULL),
(UUID(), @torneo3_id, 3, 3, 1000, NULL),
(UUID(), @torneo3_id, 4, 10, 300, NULL);

-- ===========================================
-- TORNEO 4: Torneo Completado (Histórico)
-- ===========================================
SET @torneo4_id = UUID();

INSERT INTO tournaments (
    id, name, description, tournament_type, format,
    entry_fee, prize_pool, min_rank_required, max_participants,
    current_participants, status, start_time, end_time, registration_start, registration_end
) VALUES (
    @torneo4_id,
    'Copa de Año Nuevo 2024',
    'El gran torneo de inicio de año. ¡Felicitaciones a los ganadores!',
    'special',
    'bracket',
    0,
    15000,
    'Bronce 3',
    64,
    64,
    'completed',
    DATE_SUB(NOW(), INTERVAL 10 DAY),
    DATE_SUB(NOW(), INTERVAL 3 DAY),
    DATE_SUB(NOW(), INTERVAL 15 DAY),
    DATE_SUB(NOW(), INTERVAL 11 DAY)
);

INSERT INTO tournament_prizes (id, tournament_id, position_from, position_to, token_reward, special_reward) VALUES
(UUID(), @torneo4_id, 1, 1, 7500, 'Corona de Año Nuevo'),
(UUID(), @torneo4_id, 2, 2, 4000, NULL),
(UUID(), @torneo4_id, 3, 3, 2000, NULL),
(UUID(), @torneo4_id, 4, 8, 750, NULL),
(UUID(), @torneo4_id, 9, 16, 300, NULL);

-- ===========================================
-- TORNEO 5: Próximo Torneo (Upcoming)
-- ===========================================
SET @torneo5_id = UUID();

INSERT INTO tournaments (
    id, name, description, tournament_type, format,
    entry_fee, prize_pool, min_rank_required, max_participants,
    status, start_time, end_time, registration_start, registration_end
) VALUES (
    @torneo5_id,
    'Torneo Relámpago',
    'Torneo rápido de 1 hora. ¡Acción intensa y premios instantáneos!',
    'special',
    'league',
    25,
    2500,
    'Bronce 2',
    20,
    'upcoming',
    DATE_ADD(NOW(), INTERVAL 7 DAY),
    DATE_ADD(NOW(), INTERVAL 7 DAY) + INTERVAL 1 HOUR,
    DATE_ADD(NOW(), INTERVAL 5 DAY),
    DATE_ADD(NOW(), INTERVAL 7 DAY) - INTERVAL 1 HOUR
);

INSERT INTO tournament_prizes (id, tournament_id, position_from, position_to, token_reward, special_reward) VALUES
(UUID(), @torneo5_id, 1, 1, 1500, 'Badge Velocista'),
(UUID(), @torneo5_id, 2, 2, 700, NULL),
(UUID(), @torneo5_id, 3, 3, 300, NULL);

-- Verificar torneos creados
SELECT id, name, status, current_participants, max_participants, prize_pool
FROM tournaments
ORDER BY created_at DESC;
