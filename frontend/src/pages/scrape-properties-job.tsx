import JobTable from '../components/SWRTable'
import Layout from '../components/layout'
import JobForm from '../components/scrapePropertiesJobForm'
import { useScrapePropertiesJobs, createScrapePropertiesJob } from '../models/repos/scrapePropertiesJob'
import { ScrapePropertiesJob } from '../models/types/scrapePropertiesJob'
import Badge from 'react-bootstrap/Badge';
import { useState } from 'react';
import { Button } from 'react-bootstrap'

function decorateTable(k: string, v: ScrapePropertiesJob[keyof ScrapePropertiesJob]) {
  if (k === "State") {
    return (
      <Badge bg={v === "waiting" ? "secondary" : v === "running" ? "primary" : "success"}>{v?.toString()}</Badge>
    )
  }
  if (v instanceof Date) {
    return <span>{v?.toLocaleString()}</span>
  }
  return <span>{v?.toString()}</span>
}

export default function ScrapePropertiesJobPage() {
  const [updatedJobs, setUpdatedJobs] = useState<ScrapePropertiesJob[]>([])

  return (
    <Layout home>
      <main>
        <h1 className='py-3'>物件スクレイピング ジョブ管理</h1>
      </main>

      <h2 className='pt-2'>作成</h2>
      <div className='pt-2 pb-4'>
        <JobForm onSubmited={(j) => {
          if (j) setUpdatedJobs([j])
        }}></JobForm>
      </div>
      
      <h2 className='pt-2'>一覧</h2>
      <div className='py-1
      '>
        <JobTable
          type={ScrapePropertiesJob}
          useSWRValue={useScrapePropertiesJobs}
          ignoreProperties={["DeletedAt", "Url", "Message"]}
          decorate={(k, v) => {
            if (v === undefined) return "null"
            if (k === "UpdatedAt") {
              return v.toLocaleString()
            }
            return v.toString()
          }}
          fullDecorate={decorateTable}
          maxLength={20}
          updatedJobs={updatedJobs}
          addtionalColumns={[{name: "操作", makeElement: (v) => {
            return (
              <Button
                variant="secondary"
                size="sm"
                onClick={() => {createScrapePropertiesJob(v.Url, v.Tag === undefined ? "" : v.Tag, "")}}
              >再実行</Button>
            )
          }}]}
        />
      </div>
    </Layout>
  )
}
