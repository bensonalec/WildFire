-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema outreach_school
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `outreach_school` DEFAULT CHARACTER SET utf8 ;
USE `outreach_school` ;

-- -----------------------------------------------------
-- Table `outreach_school`.`organizations`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`organizations` (
  `org_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(50) NOT NULL,
  PRIMARY KEY (`org_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `outreach_school`.`point_of_contacts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`point_of_contacts` (
  `poc_id` INT NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(50) NOT NULL,
  `last_name` VARCHAR(50) NOT NULL,
  PRIMARY KEY (`poc_id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`organizations_has_point_of_contacts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`organizations_has_point_of_contacts` (
  `org_id` INT NOT NULL,
  `poc_id` INT NOT NULL,
  PRIMARY KEY (`org_id`, `poc_id`),
  INDEX `fk_organisations_has_point_of_contacts_point_of_contacts_idx` (`poc_id` ASC),
  INDEX `fk_organisations_has_point_of_contacts_organisations_idx` (`org_id` ASC),
  CONSTRAINT `fk_organisations_has_point_of_contacts_organisations`
    FOREIGN KEY (`org_id`)
    REFERENCES `outreach_school`.`organizations` (`org_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT `fk_organisations_has_point_of_contacts_point_of_contacts`
    FOREIGN KEY (`poc_id`)
    REFERENCES `outreach_school`.`point_of_contacts` (`poc_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB;

USE `outreach_school` ;

-- -----------------------------------------------------
-- Table `outreach_school`.`events`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`events` (
  `event_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `date` DATETIME NULL DEFAULT NULL,
  `org_id` INT NOT NULL,
  PRIMARY KEY (`event_id`),
  INDEX `fk_events_organisations1_idx` (`org_id` ASC),
  CONSTRAINT `fk_events_organisations`
    FOREIGN KEY (`org_id`)
    REFERENCES `outreach_school`.`organizations` (`org_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`event_has_point_of_contacts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`event_has_point_of_contacts` (
  `event_id` INT NOT NULL,
  `poc_id` INT NOT NULL,
  PRIMARY KEY (`event_id`, `poc_id`),
  INDEX `fk_events_has_point_of_contacts_point_of_contacts_idx` (`poc_id` ASC),
  INDEX `fk_events_has_point_of_contacts_events_idx` (`event_id` ASC),
  CONSTRAINT `fk_events_has_point_of_contacts_events`
    FOREIGN KEY (`event_id`)
    REFERENCES `outreach_school`.`events` (`event_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT `fk_events_has_point_of_contacts_point_of_contacts`
    FOREIGN KEY (`poc_id`)
    REFERENCES `outreach_school`.`point_of_contacts` (`poc_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`schools`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`schools` (
  `school_id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NOT NULL,
  `city` VARCHAR(255) NOT NULL,
  `student_population` INT NULL DEFAULT NULL,
  PRIMARY KEY (`school_id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`school_has_events`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`school_has_events` (
  `school_id` INT NOT NULL,
  `event_id` INT NOT NULL,
  PRIMARY KEY (`school_id`, `event_id`),
  INDEX `fk_schools_has_events_events_idx` (`event_id` ASC),
  INDEX `fk_schools_has_events_schools_idx` (`school_id` ASC),
  CONSTRAINT `fk_schools_has_events_events`
    FOREIGN KEY (`event_id`)
    REFERENCES `outreach_school`.`events` (`event_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT `fk_schools_has_events_schools`
    FOREIGN KEY (`school_id`)
    REFERENCES `outreach_school`.`schools` (`school_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`school_has_point_of_contacts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`school_has_point_of_contacts` (
  `school_id` INT NOT NULL,
  `poc_id` INT NOT NULL,
  PRIMARY KEY (`school_id`, `poc_id`),
  INDEX `fk_schools_has_point_of_contacts_schools_idx` (`school_id` ASC),
  INDEX `fk_schools_has_point_of_contacts_point_of_contacts_idx` (`poc_id` ASC),
  CONSTRAINT `fk_schools_has_point_of_contacts_point_of_contacts`
    FOREIGN KEY (`poc_id`)
    REFERENCES `outreach_school`.`point_of_contacts` (`poc_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_schools_has_point_of_contacts_schools`
    FOREIGN KEY (`school_id`)
    REFERENCES `outreach_school`.`schools` (`school_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`students`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`students` (
  `student_id` INT NOT NULL AUTO_INCREMENT,
  `school_id` INT NOT NULL,
  `first_name` VARCHAR(50) NOT NULL,
  `last_name` VARCHAR(50) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `phone` VARCHAR(15) NULL,
  `grade` INT(2) NOT NULL,
  PRIMARY KEY (`student_id`),
  INDEX `fk_students_schools_idx` (`school_id` ASC),
  CONSTRAINT `fk_students_schools`
    FOREIGN KEY (`school_id`)
    REFERENCES `outreach_school`.`schools` (`school_id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


-- -----------------------------------------------------
-- Table `outreach_school`.`students_has_organizations`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `outreach_school`.`students_has_organizations` (
  `students_student_id` INT NOT NULL,
  `organizations_org_id` INT NOT NULL,
  `events_event_id` INT NOT NULL,
  PRIMARY KEY (`students_student_id`, `organizations_org_id`, `events_event_id`),
  INDEX `fk_students_has_organizations_organizations_idx` (`organizations_org_id` ASC),
  INDEX `fk_students_has_organizations_students_idx` (`students_student_id` ASC),
  INDEX `fk_students_has_organizations_events_idx` (`events_event_id` ASC),
  CONSTRAINT `fk_students_has_organizations_students`
    FOREIGN KEY (`students_student_id`)
    REFERENCES `outreach_school`.`students` (`student_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT `fk_students_has_organizations_organizations`
    FOREIGN KEY (`organizations_org_id`)
    REFERENCES `outreach_school`.`organizations` (`org_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  CONSTRAINT `fk_students_has_organizations_events`
    FOREIGN KEY (`events_event_id`)
    REFERENCES `outreach_school`.`events` (`event_id`)
    ON DELETE RESTRICT
    ON UPDATE CASCADE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
