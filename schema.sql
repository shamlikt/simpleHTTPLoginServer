
CREATE DATABASE apiserver;
CREATE TABLE `users` (
       `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
      `name` varchar(100) NOT NULL UNIQUE,
      `password` varchar(100) NOT NULL,
      `email` varchar(50) NOT NULL,
      `age` int(6) NOT NULL,
      `salary` int(100) NOT NULL,
      PRIMARY KEY (`id`)
      ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
