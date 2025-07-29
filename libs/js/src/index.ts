export interface QuartzflakeOptions {
    /**
     * The reference point in time. Default to January 1st, 1970. The maximum size must be 64 bits and cannot be negative.
     */
    epoch?: number,
    /**
     * The machine id for generation. The maximum size is 14 bits.
     */
    machineId?: number,
}

const DEFAULT_OPTIONS: QuartzflakeOptions = {
    epoch: 0,
    machineId: 0,
};

export class Quartzflake {
    private options: QuartzflakeOptions;
    private lastTimeStamp: number;
    private sequence: number;

    /**
     * Intialize the Quartzflake ID generator
     * @param options generation options
     */
    public constructor(options?: QuartzflakeOptions) {
        this.options = {
            ...DEFAULT_OPTIONS,
            ...options,
        };
        this.lastTimeStamp = Date.now();
        this.sequence = 0;

        if (this.options.machineId as number >= 16384) {
            throw new Error("machineId cannot be larger than 2^14");
        }

        if (this.options.epoch as number < 0) {
            throw new Error("epoch cannot be negative");
        }
    }

    /**
     * Generate Quartzflake ID
     * @returns a 96-bit bigint represented in decimal
     */
    public nextId(): bigint {
        const now: number = Date.now();
        if (now !== this.lastTimeStamp) {
            this.sequence = 0;
            this.lastTimeStamp = now;
        } else if (this.sequence >= 262144) {
            throw new Error("sequence overflow within the same millisecond");
        }

        const timestamp: bigint = BigInt(now - (this.options.epoch as number)) & ((1n << 63n) - 1n);
        const machineId: bigint = BigInt(this.options.machineId as number) & ((1n << 14n) - 1n);
        const sequence: bigint = BigInt(this.sequence) & ((1n << 18n) - 1n);

        ++this.sequence;

        return (timestamp << 32n) | (machineId << 18n) | sequence;
    }
}