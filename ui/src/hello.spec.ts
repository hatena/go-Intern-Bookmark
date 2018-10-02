import {hello} from "./hello";

describe('hello', () => {
  it('says Hello, world!', () => {
    expect(hello()).toEqual('Hello, world!');
  })
});
