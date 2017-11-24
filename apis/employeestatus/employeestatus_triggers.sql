
DROP TRIGGER IF exists `pulsecheck`.`employeestatus_AFTER_INSERT`;
DROP TRIGGER IF exists `pulsecheck`.`employeestatus_after_update`;
delimiter //
CREATE TRIGGER `pulsecheck`.`employeestatus_AFTER_INSERT` AFTER INSERT ON `employeestatus` FOR EACH ROW
BEGIN
 INSERT INTO employeestatus_a(
id, employeeid, consider, ofac,  ofaclastsearch ,  oig,  oiglastsearch,
  sam,  samlastsearch, isactive, createdby, created,modifiedby, modified
)
   VALUES (
new.id, new.employeeid, new.consider, new.ofac,  new.ofaclastsearch ,  new.oig,  new.oiglastsearch,
  new.sam,  new.samlastsearch, new.isactive,
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
id, employeeid, consider, ofac,  ofaclastsearch ,  oig,  oiglastsearch,
  sam,  samlastsearch, isactive, createdby, created,modifiedby, modified
)
   VALUES (
new.id, new.employeeid, new.consider, new.ofac,  new.ofaclastsearch ,  new.oig,  new.oiglastsearch,
  new.sam,  new.samlastsearch, new.isactive,
new.createdby, new.created, new.modifiedby, new.modified
);
END;//
delimiter ;
