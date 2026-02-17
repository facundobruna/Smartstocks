-- Fase Educacion: Sistema de Cursos
-- Migraciones para cursos, lecciones y progreso

-- ===========================================
-- TABLA: courses (Cursos disponibles)
-- ===========================================
CREATE TABLE IF NOT EXISTS courses (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(50) DEFAULT 'BookOpen',
    category ENUM('fundamentos', 'analisis', 'estrategia', 'avanzado') NOT NULL,
    difficulty ENUM('principiante', 'intermedio', 'avanzado') NOT NULL,
    duration_minutes INT DEFAULT 30,
    points_reward INT DEFAULT 100,
    is_premium BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    order_index INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_courses_category (category),
    INDEX idx_courses_difficulty (difficulty),
    INDEX idx_courses_active (is_active),
    INDEX idx_courses_order (order_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: lessons (Lecciones de cada curso)
-- ===========================================
CREATE TABLE IF NOT EXISTS lessons (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    content_type ENUM('text', 'video', 'quiz') DEFAULT 'text',
    video_url VARCHAR(500) NULL,
    duration_minutes INT DEFAULT 5,
    order_index INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_lessons_course (course_id),
    INDEX idx_lessons_order (order_index),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: lesson_quiz_questions (Preguntas de quiz en lecciones)
-- ===========================================
CREATE TABLE IF NOT EXISTS lesson_quiz_questions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    lesson_id CHAR(36) NOT NULL,
    question_text TEXT NOT NULL,
    option_a TEXT NOT NULL,
    option_b TEXT NOT NULL,
    option_c TEXT NOT NULL,
    option_d TEXT NOT NULL,
    correct_option ENUM('A', 'B', 'C', 'D') NOT NULL,
    explanation TEXT,
    order_index INT DEFAULT 0,
    INDEX idx_quiz_lesson (lesson_id),
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: user_course_progress (Progreso del usuario en cursos)
-- ===========================================
CREATE TABLE IF NOT EXISTS user_course_progress (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    course_id CHAR(36) NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    UNIQUE KEY unique_user_course (user_id, course_id),
    INDEX idx_progress_user (user_id),
    INDEX idx_progress_course (course_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: user_lesson_progress (Progreso del usuario en lecciones)
-- ===========================================
CREATE TABLE IF NOT EXISTS user_lesson_progress (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    lesson_id CHAR(36) NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    quiz_score INT NULL,
    UNIQUE KEY unique_user_lesson (user_id, lesson_id),
    INDEX idx_lesson_progress_user (user_id),
    INDEX idx_lesson_progress_lesson (lesson_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- STORED PROCEDURE: Completar leccion
-- ===========================================
DELIMITER //
CREATE PROCEDURE complete_lesson(
    IN p_user_id CHAR(36),
    IN p_lesson_id CHAR(36),
    IN p_quiz_score INT
)
BEGIN
    DECLARE v_course_id CHAR(36);
    DECLARE v_total_lessons INT;
    DECLARE v_completed_lessons INT;
    DECLARE v_course_points INT;

    -- Obtener course_id
    SELECT course_id INTO v_course_id FROM lessons WHERE id = p_lesson_id;

    -- Marcar leccion como completada
    INSERT INTO user_lesson_progress (id, user_id, lesson_id, is_completed, completed_at, quiz_score)
    VALUES (UUID(), p_user_id, p_lesson_id, TRUE, NOW(), p_quiz_score)
    ON DUPLICATE KEY UPDATE
        is_completed = TRUE,
        completed_at = NOW(),
        quiz_score = COALESCE(p_quiz_score, quiz_score);

    -- Iniciar progreso del curso si no existe
    INSERT IGNORE INTO user_course_progress (id, user_id, course_id)
    VALUES (UUID(), p_user_id, v_course_id);

    -- Contar lecciones totales y completadas
    SELECT COUNT(*) INTO v_total_lessons
    FROM lessons WHERE course_id = v_course_id AND is_active = TRUE;

    SELECT COUNT(*) INTO v_completed_lessons
    FROM user_lesson_progress ulp
    JOIN lessons l ON ulp.lesson_id = l.id
    WHERE ulp.user_id = p_user_id AND l.course_id = v_course_id AND ulp.is_completed = TRUE;

    -- Si completo todas las lecciones, marcar curso como completado
    IF v_completed_lessons = v_total_lessons THEN
        UPDATE user_course_progress
        SET is_completed = TRUE, completed_at = NOW()
        WHERE user_id = p_user_id AND course_id = v_course_id;

        -- Dar puntos por completar el curso
        SELECT points_reward INTO v_course_points FROM courses WHERE id = v_course_id;

        UPDATE user_stats
        SET smartpoints = smartpoints + v_course_points,
            updated_at = NOW()
        WHERE user_id = p_user_id;

        CALL update_user_rank(p_user_id);
    END IF;
END//
DELIMITER ;

-- ===========================================
-- SEED: Cursos iniciales
-- ===========================================
INSERT INTO courses (id, title, description, icon, category, difficulty, duration_minutes, points_reward, is_premium, order_index) VALUES
(UUID(), 'Introduccion al Mercado de Valores', 'Aprende los conceptos basicos del mercado de valores y como funcionan las acciones. Este curso te dara las bases para entender el mundo de las inversiones.', 'TrendingUp', 'fundamentos', 'principiante', 45, 500, FALSE, 1),
(UUID(), 'Ahorro e Inversion', 'Descubre la diferencia entre ahorrar e invertir, y aprende estrategias para hacer crecer tu dinero de forma inteligente.', 'PiggyBank', 'fundamentos', 'principiante', 30, 400, FALSE, 2),
(UUID(), 'Finanzas Personales', 'Aprende a manejar tu dinero, crear presupuestos efectivos y planificar tu futuro financiero con confianza.', 'Wallet', 'fundamentos', 'principiante', 75, 600, FALSE, 3),
(UUID(), 'Analisis Tecnico Basico', 'Domina el arte de leer graficos y detectar patrones en el precio de las acciones para tomar mejores decisiones.', 'BarChart3', 'analisis', 'intermedio', 60, 800, FALSE, 4),
(UUID(), 'Gestion de Riesgo', 'Estrategias probadas para proteger tu capital y minimizar perdidas en tus inversiones.', 'Shield', 'estrategia', 'intermedio', 40, 700, FALSE, 5),
(UUID(), 'Diversificacion de Portafolio', 'Como distribuir tus inversiones de forma inteligente para reducir el riesgo y maximizar rendimientos.', 'Target', 'estrategia', 'intermedio', 50, 750, TRUE, 6);

-- Insertar lecciones para el primer curso (Introduccion al Mercado de Valores)
SET @course1_id = (SELECT id FROM courses WHERE title = 'Introduccion al Mercado de Valores' LIMIT 1);

INSERT INTO lessons (id, course_id, title, content, content_type, duration_minutes, order_index) VALUES
(UUID(), @course1_id, 'Que es el Mercado de Valores?',
'# Que es el Mercado de Valores?

El **mercado de valores** es un lugar donde se compran y venden acciones de empresas. Cuando compras una accion, te conviertes en dueno de una pequena parte de esa empresa.

## Conceptos Clave

- **Accion**: Una parte de propiedad de una empresa
- **Bolsa de valores**: El lugar donde se negocian las acciones
- **Inversor**: Persona que compra acciones esperando que suban de valor

## Por que existe?

Las empresas necesitan dinero para crecer. En lugar de pedir prestamos, pueden vender partes de su empresa (acciones) al publico. Los inversores compran estas acciones esperando que la empresa crezca y sus acciones valgan mas.

## Ejemplo practico

Imagina que tu amigo tiene una panaderia exitosa y necesita $10,000 para abrir otra sucursal. En lugar de pedirte prestado, te ofrece ser "socio" - tu pones $1,000 y a cambio eres dueno del 10% de la panaderia. Si la panaderia crece, tu 10% valdra mas!',
'text', 5, 1),

(UUID(), @course1_id, 'Como funcionan las Acciones',
'# Como funcionan las Acciones

Cuando compras una **accion**, estas comprando un pedacito de una empresa. Esto te da ciertos derechos y beneficios.

## Derechos del Accionista

1. **Votar** en decisiones importantes de la empresa
2. **Recibir dividendos** (parte de las ganancias)
3. **Vender** tu accion cuando quieras

## Por que suben o bajan las acciones?

El precio de una accion depende de la **oferta y demanda**:

- Si muchos quieren comprar -> el precio **sube**
- Si muchos quieren vender -> el precio **baja**

## Factores que afectan el precio

- Ganancias de la empresa
- Noticias del sector
- Economia general
- Confianza de los inversores

## Ejemplo

Si Apple anuncia que vendio muchos iPhones, los inversores piensan que la empresa vale mas, quieren comprar acciones, y el precio sube.',
'text', 6, 2),

(UUID(), @course1_id, 'Tipos de Inversores',
'# Tipos de Inversores

No todos los inversores son iguales. Dependiendo de tus objetivos y tolerancia al riesgo, puedes ser diferente tipo de inversor.

## Inversor Conservador

- Prefiere **bajo riesgo**
- Acepta **menores ganancias** a cambio de seguridad
- Invierte en bonos, plazos fijos, acciones estables

## Inversor Moderado

- Balance entre **riesgo y ganancia**
- Diversifica sus inversiones
- Mezcla acciones con instrumentos mas seguros

## Inversor Agresivo

- Busca **altas ganancias**
- Acepta **alto riesgo**
- Invierte en acciones de crecimiento, criptomonedas

## Cual eres tu?

Preguntate:
- Cuanto dinero puedo perder sin afectar mi vida?
- Cuanto tiempo puedo esperar para ver ganancias?
- Que tan nervioso me pongo cuando mis inversiones bajan?',
'text', 5, 3),

(UUID(), @course1_id, 'La Bolsa de Buenos Aires',
'# La Bolsa de Buenos Aires (BYMA)

En Argentina, las acciones se negocian principalmente en **BYMA** (Bolsas y Mercados Argentinos).

## Historia

- Fundada en 1854
- Una de las bolsas mas antiguas de Latinoamerica
- Hoy es totalmente electronica

## Indice Merval

El **Merval** es el indice mas importante de Argentina. Mide el rendimiento de las acciones mas negociadas:

- YPF (petroleo)
- Banco Galicia
- Pampa Energia
- Telecom Argentina
- Y otras empresas lideres

## Horarios de operacion

- Lunes a viernes
- 11:00 a 17:00 (hora Argentina)
- Cerrado feriados

## Dato curioso

Podes comprar acciones argentinas desde tu celular usando apps de brokers como:
- IOL (InvertirOnline)
- Bull Market
- Balanz',
'text', 5, 4),

(UUID(), @course1_id, 'Riesgos y Beneficios',
'# Riesgos y Beneficios de Invertir

Invertir en acciones tiene ventajas y desventajas que debes conocer antes de empezar.

## Beneficios

### 1. Potencial de crecimiento
Historicamente, las acciones han dado mejores rendimientos que el ahorro tradicional a largo plazo.

### 2. Dividendos
Algunas empresas reparten parte de sus ganancias a los accionistas periodicamente.

### 3. Liquidez
Podes vender tus acciones cualquier dia que la bolsa este abierta.

### 4. Proteccion contra inflacion
Las acciones tienden a subir con la inflacion, protegiendo tu poder adquisitivo.

## Riesgos

### 1. Volatilidad
Los precios pueden subir y bajar drasticamente en poco tiempo.

### 2. Perdida de capital
Podes perder parte o todo tu dinero invertido.

### 3. Riesgo de empresa
Si la empresa quiebra, tus acciones pueden valer cero.

## Regla de oro

**Nunca inviertas dinero que necesitas a corto plazo o que no puedas permitirte perder.**',
'text', 6, 5),

(UUID(), @course1_id, 'Tu Primera Inversion',
'# Como Hacer tu Primera Inversion

Guia paso a paso para empezar a invertir en Argentina.

## Paso 1: Elegir un Broker

Un broker es la empresa que te permite comprar y vender acciones. Opciones populares:

- **IOL (InvertirOnline)**: Facil de usar, ideal para principiantes
- **Bull Market**: Buenos costos, app amigable
- **Balanz**: Amplia variedad de instrumentos

## Paso 2: Abrir una cuenta

Necesitas:
- DNI
- Comprobante de domicilio
- CBU de tu cuenta bancaria

El proceso es 100% online y tarda 24-48 horas.

## Paso 3: Depositar dinero

Transferis pesos desde tu banco a tu cuenta del broker.

## Paso 4: Elegir que comprar

Para empezar, considera:
- **CEDEARs**: Acciones de empresas extranjeras (Apple, Google, etc.)
- **FCI**: Fondos que invierten por vos
- **Acciones locales**: Empresas argentinas

## Paso 5: Comprar!

Buscas la accion, elegis cuanto comprar, y confirmas. Listo, ya sos inversor!

## Consejo

Empieza con poco dinero mientras aprendes. No hay apuro!',
'text', 8, 6),

(UUID(), @course1_id, 'Errores Comunes del Principiante',
'# Errores Comunes que Debes Evitar

Aprender de los errores de otros te ahorrara dinero y frustracion.

## Error 1: Invertir sin entender

**Problema**: Comprar acciones solo porque alguien las recomendo.
**Solucion**: Siempre investiga antes de invertir. Si no entendes el negocio, no inviertas.

## Error 2: Poner todos los huevos en una canasta

**Problema**: Invertir todo tu dinero en una sola accion.
**Solucion**: Diversifica. Distribuye tu dinero en varias inversiones.

## Error 3: Dejarse llevar por las emociones

**Problema**: Vender en panico cuando baja, comprar euforico cuando sube.
**Solucion**: Ten un plan y siguelo. Las emociones son malas consejeras.

## Error 4: Esperar hacerse rico rapido

**Problema**: Pensar que vas a duplicar tu dinero en una semana.
**Solucion**: La inversion es un maraton, no una carrera de 100 metros.

## Error 5: No tener un fondo de emergencia

**Problema**: Invertir el dinero que necesitas para vivir.
**Solucion**: Primero ahorra 3-6 meses de gastos, luego invierte.

## Recuerda

Los mejores inversores no son los mas inteligentes, son los mas disciplinados.',
'text', 6, 7),

(UUID(), @course1_id, 'Quiz Final',
'# Quiz Final: Introduccion al Mercado de Valores

Pon a prueba lo que aprendiste en este curso!

Responde las siguientes preguntas para completar el curso y ganar tus puntos.',
'quiz', 4, 8);

-- Insertar preguntas del quiz
SET @quiz_lesson_id = (SELECT id FROM lessons WHERE course_id = @course1_id AND content_type = 'quiz' LIMIT 1);

INSERT INTO lesson_quiz_questions (id, lesson_id, question_text, option_a, option_b, option_c, option_d, correct_option, explanation, order_index) VALUES
(UUID(), @quiz_lesson_id, 'Que es una accion?', 'Un prestamo a una empresa', 'Una parte de propiedad de una empresa', 'Un tipo de moneda', 'Un contrato de trabajo', 'B', 'Una accion representa una parte de la propiedad de una empresa. Cuando compras acciones, te conviertes en socio de esa empresa.', 1),
(UUID(), @quiz_lesson_id, 'Que hace que el precio de una accion suba?', 'Cuando la empresa pierde dinero', 'Cuando mas personas quieren vender', 'Cuando mas personas quieren comprar', 'El gobierno lo decide', 'C', 'El precio de las acciones sube cuando hay mas demanda (compradores) que oferta (vendedores).', 2),
(UUID(), @quiz_lesson_id, 'Que es el Merval?', 'Un banco argentino', 'El indice principal de la bolsa argentina', 'Una criptomoneda', 'Un tipo de bono', 'B', 'El Merval es el indice que mide el rendimiento de las principales acciones argentinas.', 3),
(UUID(), @quiz_lesson_id, 'Cual es la regla de oro de la inversion?', 'Invertir todo en una sola accion', 'Seguir las recomendaciones de amigos', 'Nunca invertir dinero que necesitas a corto plazo', 'Vender cuando el precio baja', 'C', 'Nunca debes invertir dinero que puedas necesitar pronto o que no puedas permitirte perder.', 4),
(UUID(), @quiz_lesson_id, 'Que tipo de inversor acepta mayor riesgo buscando mayores ganancias?', 'Conservador', 'Moderado', 'Pasivo', 'Agresivo', 'D', 'Los inversores agresivos estan dispuestos a aceptar mayor riesgo a cambio de la posibilidad de obtener mayores ganancias.', 5);
