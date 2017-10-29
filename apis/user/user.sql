CREATE TABLE `user` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`email` varchar(160) NOT NULL,
	`firstname` varchar(80) NOT NULL,
	`middlename` varchar(80) NULL,
	`lastname` varchar(80) NOT NULL,
	`companyid`	int(11) NOT NULL, 
	`isactive` tinyint(1) NOT NULL DEFAULT '1',
	`createdby` varchar(80) NOT NULL DEFAULT 'admin',
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`modifiedby` varchar(80) NOT NULL DEFAULT 'admin',
	`modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
   UNIQUE KEY `uk_user_email` (`email`),
   FOREIGN KEY (companyid) REFERENCES company(id)
) ENGINE=InnoDB;
