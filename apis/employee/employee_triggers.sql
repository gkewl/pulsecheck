-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

DROP TRIGGER IF exists `sparq`.`employee_after_insert`;
DROP TRIGGER IF exists `sparq`.`employee_after_update`;

-- +goose StatementBegin
CREATE TRIGGER `employee_after_insert`
   AFTER INSERT ON sparq.employee FOR EACH ROW
BEGIN
   INSERT INTO employee_a(
id, companyid, firstname, middlename, lastname, dateofbirth, type, isactive,
createdby, created, modifiedby, modified
)
   VALUES (
new.id, new.companyid, new.firstname, new.middlename, new.lastname,
new.dateofbirth, new.type, new.isactive, new.createdby, new.created,
new.modifiedby, new.modified
);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER `employee_after_update`
   AFTER UPDATE ON sparq.employee FOR EACH ROW
BEGIN
   INSERT INTO employee_a(
id, companyid, firstname, middlename, lastname, dateofbirth, type, isactive,
createdby, created, modifiedby, modified
)
   VALUES (
new.id, new.companyid, new.firstname, new.middlename, new.lastname,
new.dateofbirth, new.type, new.isactive, new.createdby, new.created,
new.modifiedby, new.modified
);
END;
-- +goose StatementEnd
