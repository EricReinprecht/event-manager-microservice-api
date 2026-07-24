DO
$$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_roles
        WHERE rolname = 'event_test_user'
    ) THEN
        CREATE ROLE event_test_user LOGIN PASSWORD 'event_test_password';
    END IF;
END
$$;


SELECT 'CREATE DATABASE event_platform_test OWNER event_test_user'
WHERE NOT EXISTS (
    SELECT FROM pg_database
    WHERE datname = 'event_platform_test'
)
\gexec