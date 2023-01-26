import BForm from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import { useState } from "react";
import { createScrapePropertiesJob } from '../models/repos/scrapePropertiesJob'
import { ScrapePropertiesJob } from '../models/types/scrapePropertiesJob';

export default function Form(
  { onSubmited }: { onSubmited?: (job: ScrapePropertiesJob | undefined) => void }) {
  const [url, setUrl] = useState('')
  const [tag, setTag] = useState('')
  const [args, setArgs] = useState('')

  return (
    <BForm onSubmit={async (e) => {
      e.preventDefault()
      const job = await createScrapePropertiesJob(url, tag, args)
      if (onSubmited) onSubmited(job)
    }}>
      <BForm.Group className='mb-3'>
        <BForm.Label>URL</BForm.Label>
        <BForm.Control
          type='url'
          placeholder='各検索サイトの検索結果ページURL'
          onChange={(e) => setUrl(e.target.value)}
        />
        <BForm.Label>タグ</BForm.Label>
        <BForm.Control
          type='text'
          placeholder='ジョブを区別するためのテキスト'
          onChange={(e) => setTag(e.target.value)}
        />
        <BForm.Label>オプション</BForm.Label>
        <BForm.Control
          type='text'
          placeholder='追加の引数'
          onChange={(e) => setArgs(e.target.value)}
        />
      </BForm.Group>
      <div className='d-grid gap-2'>
        <Button variant='primary' type='submit'>
          作成
        </Button>
      </div>
    </BForm>
  )
}