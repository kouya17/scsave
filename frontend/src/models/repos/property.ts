import { fetcher } from "./fetcher"
import useSWR from 'swr'
import { deserialize } from '../types/serializer'
import { ScopedMutator } from "swr/dist/types"
import { Property } from "../types/property"

const baseUrl = process.env.NEXT_PUBLIC_BACKEND_HOST
//const baseUrl = "http://localhost:8000"

export function generateUseProperty(city: string, page?: number, order?: string, minJobId?: number, maxJobId?: number) {
  return _generateUseProperty(_buildGETUrl(city, page, false, order, minJobId, maxJobId))
}

export function usePropertiesCount(city: string, minJobId: number, maxJobId: number) {
  const { data, error } = useSWR(_buildGETUrl(city, undefined, true, undefined, minJobId, maxJobId), fetcher)
  return {
    value: data ? parseInt(data, 10) : undefined,
    isLoading: !error && !data,
    isError: error
  }
}

function propertyNameToSnake(str: string) : string {
  return str.charAt(0).toLowerCase() + str.charAt(1).toLowerCase() + str.slice(2).replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`)
}

function _buildGETUrl(city: string, page?: number, isCount?: boolean, order?: string, minJobId?: number, maxJobId?: number) {
  const countPath = isCount !== undefined && isCount ? "/count" : ""
  const pageQuery = page == undefined ? "" : "&page=" + page
  const orderQuery = order == undefined ? "" : "&order=" + propertyNameToSnake(order)
  const minJobIdQuery = minJobId !== undefined && minJobId > 0 ? "&min-job=" + minJobId : ""
  const maxJobIdQuery = maxJobId !== undefined && maxJobId > 0 ? "&max-job=" + maxJobId : ""
  const url = baseUrl + '/properties' + countPath + '?city=' + city + pageQuery + orderQuery + minJobIdQuery + maxJobIdQuery
  return url
}

function _generateUseProperty(url: string) {
  return () => {
    const { data, error } = useSWR(url, fetcher)

    return {
      values: data ? deserialize(data.Properties, Property) as unknown as Property[] : undefined,
      isLoading: !error && !data,
      isError: error
    }
  }
}

export function mutateScrapePropertiesJobs(mutate: ScopedMutator<any>) {
  mutate(baseUrl + '/properties')
}

export function countUpClickCount(p: Property) {
  const body = {
    ClickCount: (p.ClickCount ? p.ClickCount : 0) + 1
  }
  fetch(baseUrl + '/properties/' + p.ID, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body)
  })
}