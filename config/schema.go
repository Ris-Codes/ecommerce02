package config

// import "log"

// func createTable(action string) {
// 	switch action {
// 	case "admin":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS admin(
// 			id BIGSERIAL PRIMARY KEY,
// 			name VARCHAR(100) UNIQUE NOT NULL,
// 			email VARCHAR UNIQUE NOT NULL,
// 			password TEXT NOT NULL,
// 			role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'superadmin'))
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create admin table")
// 		}

// 	case "users":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users(
// 			id BIGSERIAL PRIMARY KEY,
// 			username VARCHAR(100) UNIQUE NOT NULL,
// 			password TEXT NOT NULL,
// 			email VARCHAR UNIQUE NOT NULL,
// 			phone NUMERIC(10,0) UNIQUE NOT NULL
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create users table")
// 		}

// 	case "address":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS address (
// 			id BIGSERIAL PRIMARY KEY,
// 			user_id INTEGER NOT NULL,
// 			address_line1 VARCHAR(255) NOT NULL,
// 			address_line2 VARCHAR(255) NOT NULL,
// 			city VARCHAR(100) NOT NULL,
// 			state VARCHAR(100) NOT NULL,
// 			postal_code VARCHAR(20) NOT NULL,
// 			country VARCHAR(100) NOT NULL,
// 			is_default BOOLEAN  DEFAULT false,
// 			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create address table")
// 		}

// 	case "products":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS products (
// 			id BIGSERIAL PRIMARY KEY,
// 			category_id INTEGER NOT NULL,
// 			product_name TEXT NOT NULL,
// 			description TEXT NOT NULL,
// 			price INTEGER NOT NULL,
// 			stock INTEGER NOT NULL,
// 			product_image TEXT,
// 			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create products table")
// 		}

// 	case "categories":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS categories (
// 			id BIGSERIAL PRIMARY KEY,
// 			category_name TEXT NOT NULL
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create categories table")
// 		}

// 	case "cart":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS cart (
// 			id BIGSERIAL PRIMARY KEY,
// 			user_id INTEGER NOT NULL,
// 			product_id INTEGER NOT NULL,
// 			quantity INTEGER NOT NULL DEFAULT 1,
// 			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
// 			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create cart table")
// 		}

// 	case "orders":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS orders (
// 			id BIGSERIAL PRIMARY KEY,
// 			ref_number VARCHAR(50) NOT NULL,
// 			user_id INTEGER NOT NULL,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			shipping_address_id INTEGER NOT NULL,
// 			status_id INTEGER NOT NULL,
// 			order_total INTEGER NOT NULL,
// 			payment_intent_id INTEGER NOT NULL,
// 			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
// 			FOREIGN KEY (shipping_address_id) REFERENCES address(id),
// 			FOREIGN KEY (status_id) REFERENCES order_status(id)
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create orders table")
// 		}

// 	case "order_items":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS order_items (
// 			id BIGSERIAL PRIMARY KEY,
// 			order_id INTEGER NOT NULL,
// 			product_id INTEGER NOT NULL,
// 			quantity INTEGER NOT NULL,
// 			price INTEGER NOT NULL,
// 			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
// 			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create order_items table")
// 		}

// 	case "orders_status":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS order_status (
// 			id BIGSERIAL PRIMARY KEY,
// 			status VARCHAR(50) NOT NULL
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create orders_status table")
// 		}

// 	case "payment_methods":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS payment_methods (
// 			id BIGSERIAL PRIMARY KEY,
// 			method TEXT NOT NULL
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create payment_methods table")
// 		}

// 	case "payment_status":
// 		_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS payment_status (
// 			id BIGSERIAL PRIMARY KEY,
// 			status TEXT NOT NULL
// 		)`)
// 		if err != nil {
// 			log.Fatal("failed to create payment_status table")
// 		}
// 	}
// }
