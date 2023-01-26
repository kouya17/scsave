import { JsonSerializer } from 'typescript-json-serializer';

const serializer = new JsonSerializer()

export function deserialize<T extends Object>(data: string | object, type: T) {
  return serializer.deserialize(data, type)
}