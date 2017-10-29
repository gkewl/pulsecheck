CREATE TABLE `employee` (
	`id` bigint(20) NOT NULL AUTO_INCREMENT,
	`companyid`	int(11) not null, 
	`firstname` varchar(80) NOT NULL,
	`middlename` varchar(80) NULL,
	`lastname` varchar(80) NOT NULL,
	`dateofbirth` DATE NOT NULL ,
	`type` smallint (4) NOT NULL default 1 , -- 1 - perm, 2-temp
	`isactive` tinyint(1) NOT NULL DEFAULT '1',
	`createdby` varchar(80) NOT NULL DEFAULT 'admin',
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`modifiedby` varchar(80) NOT NULL DEFAULT 'admin',
	`modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
   FOREIGN KEY (companyid) REFERENCES company(id)
) ENGINE=InnoDB;
