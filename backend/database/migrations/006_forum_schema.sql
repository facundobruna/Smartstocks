-- Fase 3: Sistema de Foro Comunitario
-- Migraciones para tablas de posts, replies y reacciones

-- ===========================================
-- TABLA: forum_posts (Publicaciones del foro)
-- ===========================================
CREATE TABLE IF NOT EXISTS forum_posts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category ENUM('general', 'inversiones', 'ahorro', 'mercados', 'cripto', 'economia', 'preguntas', 'noticias') DEFAULT 'general',
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    views INT DEFAULT 0,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_posts_user (user_id),
    INDEX idx_posts_category (category),
    INDEX idx_posts_pinned (is_pinned),
    INDEX idx_posts_created (created_at),
    INDEX idx_posts_likes (likes),
    FULLTEXT INDEX idx_posts_search (title, content),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: forum_replies (Respuestas a posts)
-- ===========================================
CREATE TABLE IF NOT EXISTS forum_replies (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    post_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    parent_reply_id CHAR(36) NULL,
    content TEXT NOT NULL,
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    is_solution BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_replies_post (post_id),
    INDEX idx_replies_user (user_id),
    INDEX idx_replies_parent (parent_reply_id),
    INDEX idx_replies_created (created_at),
    FOREIGN KEY (post_id) REFERENCES forum_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_reply_id) REFERENCES forum_replies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TABLA: forum_reactions (Likes/Dislikes)
-- ===========================================
CREATE TABLE IF NOT EXISTS forum_reactions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    post_id CHAR(36) NULL,
    reply_id CHAR(36) NULL,
    is_like BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_reactions_user (user_id),
    INDEX idx_reactions_post (post_id),
    INDEX idx_reactions_reply (reply_id),
    UNIQUE KEY unique_user_post (user_id, post_id),
    UNIQUE KEY unique_user_reply (user_id, reply_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES forum_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_id) REFERENCES forum_replies(id) ON DELETE CASCADE,
    CHECK (post_id IS NOT NULL OR reply_id IS NOT NULL)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- TRIGGER: Actualizar reply_count en posts
-- ===========================================
DELIMITER //
CREATE TRIGGER after_reply_insert
AFTER INSERT ON forum_replies
FOR EACH ROW
BEGIN
    UPDATE forum_posts
    SET reply_count = reply_count + 1
    WHERE id = NEW.post_id;
END//
DELIMITER ;

DELIMITER //
CREATE TRIGGER after_reply_delete
AFTER DELETE ON forum_replies
FOR EACH ROW
BEGIN
    UPDATE forum_posts
    SET reply_count = reply_count - 1
    WHERE id = OLD.post_id;
END//
DELIMITER ;

-- ===========================================
-- TRIGGER: Actualizar likes/dislikes en posts
-- ===========================================
DELIMITER //
CREATE TRIGGER after_reaction_insert
AFTER INSERT ON forum_reactions
FOR EACH ROW
BEGIN
    IF NEW.post_id IS NOT NULL THEN
        IF NEW.is_like THEN
            UPDATE forum_posts SET likes = likes + 1 WHERE id = NEW.post_id;
        ELSE
            UPDATE forum_posts SET dislikes = dislikes + 1 WHERE id = NEW.post_id;
        END IF;
    ELSEIF NEW.reply_id IS NOT NULL THEN
        IF NEW.is_like THEN
            UPDATE forum_replies SET likes = likes + 1 WHERE id = NEW.reply_id;
        ELSE
            UPDATE forum_replies SET dislikes = dislikes + 1 WHERE id = NEW.reply_id;
        END IF;
    END IF;
END//
DELIMITER ;

DELIMITER //
CREATE TRIGGER after_reaction_delete
AFTER DELETE ON forum_reactions
FOR EACH ROW
BEGIN
    IF OLD.post_id IS NOT NULL THEN
        IF OLD.is_like THEN
            UPDATE forum_posts SET likes = likes - 1 WHERE id = OLD.post_id;
        ELSE
            UPDATE forum_posts SET dislikes = dislikes - 1 WHERE id = OLD.post_id;
        END IF;
    ELSEIF OLD.reply_id IS NOT NULL THEN
        IF OLD.is_like THEN
            UPDATE forum_replies SET likes = likes - 1 WHERE id = OLD.reply_id;
        ELSE
            UPDATE forum_replies SET dislikes = dislikes - 1 WHERE id = OLD.reply_id;
        END IF;
    END IF;
END//
DELIMITER ;

-- ===========================================
-- STORED PROCEDURE: Incrementar vistas
-- ===========================================
DELIMITER //
CREATE PROCEDURE increment_post_views(IN p_post_id CHAR(36))
BEGIN
    UPDATE forum_posts SET views = views + 1 WHERE id = p_post_id;
END//
DELIMITER ;
