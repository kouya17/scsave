import Layout from '../components/layout'
import PropertySearchForm, { SearchOrder } from '../components/propertySearchForm'
import { useState } from "react";
import PropertyTable from '../components/propetyTable'

import Pagination from 'react-bootstrap/Pagination';
import { Dispatch, SetStateAction } from "react";
import { usePropertiesCount } from '../models/repos/property'
import { useScrapePropertiesJobs } from '../models/repos/scrapePropertiesJob';

function PropertiesPagenation({ now, max, setPage }: { now: number, max: number, setPage: Dispatch<SetStateAction<number>> }) {
  return (
    <Pagination>
      <Pagination.First onClick={() => setPage(1)} />
      <Pagination.Prev onClick={() => setPage(now - 1)} />
      {now > 1 && <Pagination.Item onClick={() => setPage(1)}>1</Pagination.Item>}
      {now > 3 && <Pagination.Ellipsis />}
      {now > 2 && <Pagination.Item onClick={() => setPage(now - 1)}>{now - 1}</Pagination.Item>}
      <Pagination.Item active>{now}</Pagination.Item>
      {now < (max - 1) && <Pagination.Item onClick={() => setPage(now + 1)}>{now + 1}</Pagination.Item>}
      {now < (max - 2) && <Pagination.Ellipsis />}
      {now < max && <Pagination.Item onClick={() => setPage(max)}>{max}</Pagination.Item>}
      <Pagination.Next onClick={() => setPage(now + 1)} />
      <Pagination.Last onClick={() => setPage(max)} />
    </Pagination>
  )
}

export default function PropertyPage() {
  const [city, setCity] = useState('')
  const [orders, setOrders] = useState<SearchOrder[]>([])
  const [minJobId, setMinJobId] = useState(-1)
  const [maxJobId, setMaxJobId] = useState(-1)
  const [page, setPage] = useState(1)
  const { value, isLoading, isError } = usePropertiesCount(city, minJobId, maxJobId)
  const useJobs = useScrapePropertiesJobs()

  return (
    <Layout>
      <main>
        <h1 className='py-3'>物件一覧</h1>
      </main>

      <div className='py-1'>
        <PropertySearchForm onSubmit={(city, orders, minJobId, maxJobId) => {
          setCity(city)
          setOrders([...orders]) // force update
          setMinJobId(minJobId)
          setMaxJobId(maxJobId)
          setPage(1)
        }}/>
      </div>
      
      <div className='py-4'>
        <PropertiesPagenation now={page} max={value ? Math.ceil(value / 10) : 0} setPage={setPage}/>
        <PropertyTable cond={{order: orders[0], wheres: [{
            property: "City", op: "LIKE", value: city
          },{
            property: "JobId", op: ">=", value: minJobId
          },{
            property: "JobId", op: "<=", value: maxJobId
          }], page: page, lastJobId: useJobs.values?.reduce((a, c) => a > c.ID ? a : c.ID, 0)}}/>
      </div>
    </Layout>
  )
}
