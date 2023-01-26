import JobTable from '../components/SWRTable'
import { Property } from '../models/types/property'
import { countUpClickCount, generateUseProperty, usePropertiesCount } from '../models/repos/property'
import { CSSProperties } from 'react'

type OrderBy = {
  property: string
  order: string // should be "desc" | "asc"
}

type Where = {
  property: string
  op: "LIKE" | "<=" | ">="
  value: string | number
}
type Condition = {
  order?: OrderBy
  wheres: Where[]
  page: number
  lastJobId?: number
}

function Table({ cond }: { cond: {city: string, page: number, order?: string, lastJobId?: number, minJobId: number, maxJobId: number} }) {
  return (
    <JobTable
      type={Property}
      useSWRValue={generateUseProperty(cond.city, cond.page, cond.order, cond.minJobId, cond.maxJobId)}
      ignoreProperties={[
        "ID",
        "CreatedAt",
        "UpdatedAt",
        "DeletedAt",
        "Url",
        "Station",
        "CoverageRatio",
        "OtherCost",
        "Timing",
        "Rights",
        "Structure",
        "BuildCompany",
        "Reform",
        "LandKind",
        "OtherRestriction",
        "OtherNotice",
      ]}
      decorate={(k, v) => {
        if (v === undefined) return "null"
        if (v instanceof Date) {
          return v.toLocaleString()
        }
        return v.toString()
      }}
      maxLength={20}
      rowDecorate={(v, i) => {
        var ret: CSSProperties = {}
        ret.cursor = 'pointer'
        if (v.ClickCount > 0) ret.backgroundColor = 'khaki'
        if (cond.lastJobId && v.JobId === cond.lastJobId) ret.borderColor = 'red'
        return ret
      }}
      onClick={(v, i) => {
        return () => {
          countUpClickCount(v)
          open(v.Url)
        }
      }}
    />
  )
}

export default function PropertyTable({ cond }: { cond: Condition }) {
  var order = ""
  if (cond.order) order = cond.order.property + " " + cond.order.order
  const city = cond.wheres.filter(v => v.property === "City")[0].value.toString()
  const minJobId = cond.wheres.filter(v => v.property === "JobId" && v.op === ">=")[0].value as number
  const maxJobId = cond.wheres.filter(v => v.property === "JobId" && v.op === "<=")[0].value as number
  
  return (
    <div>
      <Table cond={{
        city: city,
        page: cond.page,
        order: order,
        lastJobId: cond.lastJobId,
        minJobId: minJobId,
        maxJobId: maxJobId
      }}/>
    </div>
  )
}