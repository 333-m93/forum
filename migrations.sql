-- Créer les tables pour le forum

CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
  id VARCHAR(255) PRIMARY KEY,
  user_id INT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_expires (expires_at)
);

CREATE TABLE IF NOT EXISTS categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
  id INT AUTO_INCREMENT PRIMARY KEY,
  category_id INT NOT NULL,
  user_id INT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_category_time (category_id, created_at DESC)
);

-- Insérer les catégories par défaut
INSERT IGNORE INTO categories (name, description) VALUES
('Chat général', 'Discussions générales sur tous sujets'),
('MMA', 'Discussions et ressources sur MMA'),
('Boxe Anglaise', 'Discussions et ressources sur Boxe Anglaise'),
('Muay Thai', 'Discussions et ressources sur Muay Thai'),
('Jujitsu Brésilien', 'Discussions et ressources sur Jujitsu Brésilien'),
('Grappling', 'Discussions et ressources sur Grappling'),
('Autres sports de combat', 'Discussions sur autres sports de combat');
