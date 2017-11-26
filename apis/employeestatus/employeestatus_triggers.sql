
DROP TRIGGER IF exists `pulsecheck`.`employeestatus_AFTER_INSERT`;
DROP TRIGGER IF exists `pulsecheck`.`employeestatus_after_update`;
delimiter //
CREATE TRIGGER `pulsecheck`.`employeestatus_AFTER_INSERT` AFTER INSERT ON `employeestatus` FOR EACH ROW
BEGIN
 INSERT INTO employeestatus_a(
id, employeeid, consider, ofac,  ofaclastsearch , ofacreference,  oig,  oiglastsearch,oigreference,
  sam,  samlastsearch,samreference, isactive, createdby, created,modifiedby, modified
)
   VALUES (
new.id, new.employeeid, new.consider, new.ofac,  new.ofaclastsearch , new.ofacreference, new.oig,  new.oiglastsearch,new.oigreference, 
  new.sam,  new.samlastsearch, new.samreference,new.isactive,
new.createdby, new.created, new.modifiedby, new.modified
);
END;//
delimiter ;


-- +goose StatementBegin

delimiter //
CREATE  TRIGGER `pulsecheck`.`employeestatus_after_update`
   AFTER UPDATE ON `employeestatus` FOR EACH ROW
BEGIN
   INSERT INTO employeestatus_a(
id, employeeid, consider, ofac,  ofaclastsearch , ofacreference,  oig,  oiglastsearch,oigreference,
  sam,  samlastsearch,samreference, isactive, createdby, created,modifiedby, modified
)
   VALUES (
new.id, new.employeeid, new.consider, new.ofac,  new.ofaclastsearch , new.ofacreference, new.oig,  new.oiglastsearch,new.oigreference, 
  new.sam,  new.samlastsearch, new.samreference,new.isactive,
new.createdby, new.created, new.modifiedby, new.modified
);
END;//
delimiter ;
