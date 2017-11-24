CREATE TABLE `employeestatus` (
	`id` bigint(20) NOT NULL AUTO_INCREMENT,
	`employeeid` bigint(20) NOT NULL,
  `consider` tinyint(1) GENERATED ALWAYS AS ((`oig` or `sam` or `ofac`)) VIRTUAL,
  `ofac` tinyint(1) NOT NULL DEFAULT '0',
  `ofaclastsearch` timestamp NULL DEFAULT NULL,
  `oig` tinyint(1) NOT NULL DEFAULT '0',
  `oiglastsearch` timestamp NULL DEFAULT NULL,
  `sam` tinyint(1) NOT NULL DEFAULT '0',
  `samlastsearch` timestamp NULL DEFAULT NULL,
	`isactive` tinyint(1) NOT NULL DEFAULT '1',
	`createdby` varchar(80) NOT NULL DEFAULT 'admin',
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`modifiedby` varchar(80) NOT NULL DEFAULT 'admin',
	`modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;


CREATE  TABLE `employeestatus_a` (
  `auditid` bigint(20) NOT NULL AUTO_INCREMENT,
  `id` bigint(20) NOT NULL ,
  `employeeid` bigint(20) NOT NULL,
  `consider` tinyint(1) NOT NULL ,
  `ofac` tinyint(1) NOT NULL DEFAULT '0',
  `ofaclastsearch` timestamp NULL DEFAULT NULL,
  `oig` tinyint(1) NOT NULL DEFAULT '0',
  `oiglastsearch` timestamp NULL DEFAULT NULL,
  `sam` tinyint(1) NOT NULL DEFAULT '0',
  `samlastsearch` timestamp NULL DEFAULT NULL,
  `isactive` tinyint(1) NOT NULL DEFAULT '1',
  `createdby` varchar(80) NOT NULL DEFAULT 'admin',
  `created` timestamp NOT NULL ,
  `modifiedby` varchar(80) NOT NULL ,
  `modified` timestamp NOT NULL  ,
  PRIMARY KEY (`auditid`),
  KEY `fk_employeestatus_a_employeeid_idx` (`employeeid`),
  KEY `idx_employeestatus_a_consider` (`consider`),
   KEY `idx_employeestatus_a_id` (`id`),
  CONSTRAINT `fk_employeestatus_a_employeeid` FOREIGN KEY (`employeeid`) REFERENCES `employee` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB ;
