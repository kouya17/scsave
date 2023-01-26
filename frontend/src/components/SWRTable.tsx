import { CSSProperties, MouseEventHandler, useEffect } from 'react';
import BTable from 'react-bootstrap/Table';
import { useSWRConfig } from 'swr';
import { mutateScrapePropertiesJobs } from '../models/repos/scrapePropertiesJob';
import { ScrapePropertiesJob } from '../models/types/scrapePropertiesJob';
import { clip } from '../utils/string';

export default function Table<T>(
  { type, useSWRValue, ignoreProperties, decorate, fullDecorate, maxLength, updatedJobs, rowDecorate, onClick, addtionalColumns }: {
    type: {new(): T},
    useSWRValue: () => {
      values: T[] | undefined;
      isLoading: boolean;
      isError: any;
    },
    ignoreProperties?: string[],
    decorate?: (key: string, value: T[keyof T]) => string,
    fullDecorate?: (key: string, value: T[keyof T]) => JSX.Element,
    maxLength?: number,
    updatedJobs?: ScrapePropertiesJob[],
    rowDecorate?: (value: T, i: number) => CSSProperties,
    onClick?: (value: T, i: number) => MouseEventHandler<HTMLTableRowElement>,
    addtionalColumns?: {name: string, makeElement: (v: T) => JSX.Element}[]
  }
) {
  const {values, isLoading, isError} = useSWRValue()
  const {mutate} = useSWRConfig()
  useEffect(() => {
    mutateScrapePropertiesJobs(mutate)
  }, [updatedJobs])

  var maxClipLength = 10
  if (maxLength) maxClipLength = maxLength
  
  if (isLoading) return <div>loading</div>
  if (isError) return <div>error</div>
  return (
    <BTable striped bordered hover>
      <thead>
        <tr>
          {Object.getOwnPropertyNames(new type())
            .filter(e => !ignoreProperties?.includes(e))
            .map(e => {
              return <th key={e}>{e}</th>
            }
          )}
          {addtionalColumns?.map(e => {
            return <th key={e.name}>{e.name}</th>
          })}
        </tr>
      </thead>
      <tbody>
        {(values)?.map((value, i) => {
          return (
            <tr key={i} onClick={onClick && onClick(value, i)} style={rowDecorate && rowDecorate(value, i)}>
              {Object.getOwnPropertyNames(new type()).filter(e => !ignoreProperties?.includes(e)).map((key, j) => {
                return (
                  <th key={j}>
                    {fullDecorate
                      ? fullDecorate(key, value[key as keyof T])
                      : decorate
                        ? clip(decorate(key, value[key as keyof T]).toString(), maxClipLength)
                        : clip(value[key as keyof T]?.toString(), maxClipLength)
                    }
                  </th>
                )
              })}
              {addtionalColumns?.map(e=> {
                return <th>{e.makeElement(value)}</th>
              })}
            </tr>
          )
        })}
      </tbody>
    </BTable>
  )
}