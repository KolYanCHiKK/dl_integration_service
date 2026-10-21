-- Создание типа: auth_method_type
DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type t
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE t.typname = 'auth_method_type' AND n.nspname = 'edp'
	) THEN
        CREATE TYPE edp.auth_method_type AS ENUM
            ('client_secret_basic', 'client_secret_post', 'none', 'private_key_jwt');
    END IF;
END $$;

-- Создание типа: client_status_type
DO $$
BEGIN
	IF NOT EXISTS (
	    SELECT 1 FROM pg_type t
	    JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE t.typname = 'client_status_type' AND n.nspname = 'edp'
	) THEN
        CREATE TYPE edp.client_status_type AS ENUM
            ('active', 'blocked', 'deleted');
    END IF;
END $$;

-- Создание типа: grant_type_enum
DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type t
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE t.typname = 'grant_type_enum' AND n.nspname = 'edp'
	) THEN
        CREATE TYPE edp.grant_type_enum AS ENUM
            ('authorization_code', 'refresh_token', 'client_credentials');
    END IF;
END $$;

-- Создание типа: responce_type_enum
DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type t
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE t.typname = 'responce_type_enum' AND n.nspname = 'edp'
	) THEN
        CREATE TYPE edp.grant_type_enum AS ENUM
        ('code', 'id_token', 'token');
    END IF;
END $$;