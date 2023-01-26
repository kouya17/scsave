import { Form, InputGroup, Button } from "react-bootstrap";
import { Property } from "../models/types/property";
import { useState } from "react";

export type SearchOrder = {
  property: string,
  order: string
}

export default function PropertySearchForm({ onSubmit }: {
  onSubmit: (city: string, orders: SearchOrder[], minJobId: number, maxJobId: number) => void
}) {
  const [city, setCity] = useState('')
  const [orders, setOrders] = useState<SearchOrder[]>([{property: "ID", order: "desc"}])
  const [minJobId, setMinJobId] = useState(-1)
  const [maxJobId, setMaxJobId] = useState(-1)

  return (
    <Form onSubmit={(e) => {
      e.preventDefault()
      onSubmit(city, orders, minJobId, maxJobId)
    }}>
      <Form.Group>
        <Form.Label>地名</Form.Label>
        <Form.Control type='text' onChange={(e) => {
          setCity(e.target.value)
        }} />
        <Form.Label>表示順</Form.Label>
        <InputGroup>
          <Form.Select onChange={(e) => {
            const newOrders = orders
            if (!newOrders[0]) newOrders.push({property: "", order: ""})
            newOrders[0].property = e.target.value
            setOrders(newOrders)
          }}>
            {Object.getOwnPropertyNames(new Property()).map((e) =>{
              return <option key={e}>{e}</option>
            })}
          </Form.Select>
          <Form.Select onChange={(e) => {
            const newOrders = orders
            if (!newOrders[0]) newOrders.push({property: "", order: ""})
            newOrders[0].order = e.target.value
            setOrders(newOrders)
          }}>
            <option>desc</option>
            <option>asc</option>
          </Form.Select>
        </InputGroup>
        <Form.Label>JobId</Form.Label>
        <InputGroup>
          <Form.Control type="number" placeholder="min" onChange={(e) => setMinJobId(parseInt(e.target.value))}/>
          <Form.Control type="number" placeholder="max" onChange={(e) => setMaxJobId(parseInt(e.target.value))}/>
        </InputGroup>
      </Form.Group>
      <div className='mt-3 d-grid gap-2'>
        <Button variant='primary' type='submit'>
          更新
        </Button>
      </div>
    </Form>
  )
}