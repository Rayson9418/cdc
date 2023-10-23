db.getSiblingDB("demo").createCollection("demo1");
db.getSiblingDB("demo").createCollection("demo2");
db.getSiblingDB("record").createCollection("position");

// Inserting data into the demo1 collection
db.getSiblingDB("demo").getCollection("demo1").insertMany([
    { event_time: new Date("2023-10-13T08:00:00Z"), event_name: "Event 1", event_desc: "Description of event 1" },
    { event_time: new Date("2023-10-13T09:00:00Z"), event_name: "Event 2", event_desc: "Description of event 2" },
    { event_time: new Date("2023-10-13T10:00:00Z"), event_name: "Event 3", event_desc: "Description of event 3" },
    { event_time: new Date("2023-10-13T11:00:00Z"), event_name: "Event 4", event_desc: "Description of event 4" },
    { event_time: new Date("2023-10-13T12:00:00Z"), event_name: "Event 5", event_desc: "Description of event 5" }
]);

// Inserting data into the demo2 collection
db.getSiblingDB("demo").getCollection("demo2").insertMany([
    { name: "John Doe", age: 30, email: "john@example.com", address: "123 Street, City, Country" },
    { name: "Jane Smith", age: 25, email: "jane@example.com", address: "456 Avenue, Town, Country" },
    { name: "Michael Johnson", age: 35, email: "michael@example.com", address: "789 Road, Village, Country" },
    { name: "Emily Williams", age: 28, email: "emily@example.com", address: "1011 Lane, County, Country" },
    { name: "Robert Brown", age: 40, email: "robert@example.com", address: "1213 Drive, State, Country" }
]);