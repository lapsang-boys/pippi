import { Buffer } from "buffer";

export class Selection {
  start: number;
  end: number;
  content: Buffer;

  constructor(
    start?,
    end?,
    content?
  ) {
    this.start = start || 0;
    this.end = end || 0;
    this.content = content || Buffer.from([]);
  }
}
