CREATE USER 'hba_user'@'%' IDENTIFIED BY 'hbaPass123';

GRANT ALL PRIVILEGES ON `hba_db`.* TO 'hba_user'@'%';

FLUSH PRIVILEGES;

CREATE DATABASE IF NOT EXISTS `hba_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `hba_db`;

CREATE TABLE `zombies` (
                           `SessionId` varchar(10) NOT NULL,
                           `RemoteAddr` varchar(15) NOT NULL,
                           `RemotePort` varchar(6) NOT NULL,
                           `UserName` varchar(20) NOT NULL,
                           `Country` varchar(2) NOT NULL,
                           `FirstConnTime` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
                           `LastConnTime` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

