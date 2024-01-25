#! /bin/bash

set -e

createdb federation_sample

psql federation_sample -e <<-EOS
  CREATE TABLE people (id character varying, name character varying, domain character varying);
  INSERT INTO people (id, name, domain) VALUES 
    ('GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG', 'bob', 'diamcircle.org'),
    ('GCYMGWPZ6NC2U7SO6SMXOP5ZLXOEC5SYPKITDMVEONLCHFSCCQR2J4S3', 'alice', 'diamcircle.org');
EOS