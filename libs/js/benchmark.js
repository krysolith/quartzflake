const { Quartzflake } = require("./dist");
const { Bench } = require("tinybench");

const quartz = new Quartzflake();

const bench = new Bench({
    name: "Quartzflake IDs per second (1s)",
    time: 1000,
});

bench.add("Quartzflake#nextId", () => {
    quartz.nextId();
});

bench.runSync();

console.log(bench.name);
console.log(bench.table());