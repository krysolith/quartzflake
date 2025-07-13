# Quartzflake
Quartzflake is a distributed unique ID generator based on [Twitter's Snowflak](https://blog.x.com/engineering/en_us/a/2010/announcing-snowflake) and [Sonyflake](https://github.com/sony/sonyflake).

Quartzflake focuses on the lifetime of the ID, which is 576460 years, and the ID is 96 bits long. The ID is composed of the following parts:
- Timestamp (64 bits): The number of milliseconds since the Unix epoch (January 1, 1970).
- Machine ID (14 bits): A unique identifier for the machine generating the ID.
- Sequence (18 bits): A sequence number that is incremented for each ID generated within the same millisecond.

And all of them consisted in 2 parts ``High`` and ``Low`` or ``Timestamp`` and ``Metadata``.

As a result, Quartzflake has:
- A very long lifetime of 576460 years.
- A very high throughput of 2^18 IDs per millisecond.
