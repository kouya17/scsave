import { fetcher } from "./fetcher"
import useSWR from 'swr'
import { deserialize } from '../types/serializer'
import { ScrapePropertiesJob } from "../types/scrapePropertiesJob"
import { ScopedMutator } from "swr/dist/types"

const baseUrl = process.env.NEXT_PUBLIC_BACKEND_HOST
//const baseUrl = "http://localhost:8000"

export function useScrapePropertiesJobs() {
  const {data, error} = useSWR(baseUrl + '/scrape-properties-jobs', fetcher)

  return {
    values: data ? deserialize(data, ScrapePropertiesJob) as unknown as ScrapePropertiesJob[] : undefined,
    isLoading: !error && !data,
    isError: error
  }
}

export function mutateScrapePropertiesJobs(mutate: ScopedMutator<any>) {
  mutate(baseUrl + '/scrape-properties-jobs')
}

export async function createScrapePropertiesJob(url: string, tag: string, args: string) {
  const res = await fetch(baseUrl + '/scrape-properties-jobs', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({url: url, tag: tag, args: args})
  })
  const json = await res.json()
  return json ? deserialize(json, ScrapePropertiesJob) as unknown as ScrapePropertiesJob : undefined
}
