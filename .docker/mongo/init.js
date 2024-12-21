db = new Mongo().getDB("saga");
db.createCollection("workflows");