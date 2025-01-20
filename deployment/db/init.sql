// deployment/db/init.sql
CREATE TABLE IF NOT EXISTS individuals (
    id VARCHAR(255) PRIMARY KEY,
    href VARCHAR(255),
    title VARCHAR(50),
    given_name VARCHAR(255) NOT NULL,
    family_name VARCHAR(255),
    marital_status VARCHAR(1),
    gender VARCHAR(1),
    name_type VARCHAR(50),
    nationality VARCHAR(3),
    creation_date TIMESTAMP NOT NULL,
    modification_date TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contact_media (
    id VARCHAR(255) PRIMARY KEY,
    individual_id VARCHAR(255) REFERENCES individuals(id),
    type VARCHAR(50),
    medium_type VARCHAR(50),
    preferred BOOLEAN,
    phone_number VARCHAR(50),
    street1 VARCHAR(255),
    street2 VARCHAR(255),
    city VARCHAR(255),
    state_or_province VARCHAR(255),
    country VARCHAR(255),
    post_code VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS external_references (
    id VARCHAR(255) PRIMARY KEY,
    individual_id VARCHAR(255) REFERENCES individuals(id),
    name VARCHAR(255),
    external_identifier_type VARCHAR(50),
    type VARCHAR(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS individual_identifications (
    id VARCHAR(255) PRIMARY KEY,
    individual_id VARCHAR(255) REFERENCES individuals(id),
    identification_type VARCHAR(50),
    identification_id VARCHAR(255),
    valid_for_end TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS party_characteristics (
    id VARCHAR(255) PRIMARY KEY,
    individual_id VARCHAR(255) REFERENCES individuals(id),
    name VARCHAR(255),
    value VARCHAR(255),
    value_type VARCHAR(50),
    type VARCHAR(50),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_individuals_given_name ON individuals(given_name);
CREATE INDEX idx_individuals_family_name ON individuals(family_name);
CREATE INDEX idx_contact_media_individual_id ON contact_media(individual_id);
CREATE INDEX idx_external_references_individual_id ON external_references(individual_id);
CREATE INDEX idx_individual_identifications_individual_id ON individual_identifications(individual_id);
CREATE INDEX idx_party_characteristics_individual_id ON party_characteristics(individual_id);