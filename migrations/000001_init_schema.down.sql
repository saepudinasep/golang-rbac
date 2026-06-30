-- 1. Hapus tabel users terlebih dahulu karena memiliki Foreign Key yang bergantung pada tabel roles
DROP TABLE IF EXISTS users;

-- 2. Setelah users terhapus, baru hapus tabel roles
DROP TABLE IF EXISTS roles;