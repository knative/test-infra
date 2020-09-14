/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

CREATE TABLE Clusters (
  ID int NOT NULL AUTO_INCREMENT,
  ProjectID varchar(1023) NOT NULL,
  Status varchar(1023) DEFAULT 'WIP',
  Zone varchar(1023) NOT NULL,
  Nodes int NOT NULL,
  NodeType varchar(1023) NOT NULL,
  PRIMARY KEY (ID)
);

CREATE TABLE Requests (
  ID int NOT NULL AUTO_INCREMENT,
  AccessToken varchar(1023) NOT NULL,
  RequestTime timestamp,
  Zone varchar(1023) NOT NULL,
  Nodes varchar(1023) NOT NULL,
  NodeType varchar(1023) NOT NULL,
  ProwJobID varchar(1023) NOT NULL,
  ClusterID int DEFAULT 0,
  CONSTRAINT FkCluster 
  FOREIGN KEY (ClusterID) 
      REFERENCES Clusters(ID)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
  PRIMARY KEY (ID)
);
