import axios from 'axios'
import { HOST_PATH } from '../parameters'
import { getAuth } from '../../utils'

const getBody = (term, page) => {
  const currentDate = new Date()
  const dateLess30m = new Date(currentDate - (30 * 60 * 1000))

  console.log("page: ", page)
  if (!page.count)
    throw new Error("The count property must be set")

  const _from = page
    ? (page.from-1) * page.size
    : 0
  return (
    {
      "query": {
        "bool": {
          "must": [
            {
              "range": {
                "@timestamp": {
                  "lt": currentDate.toISOString(),
                  "gte": dateLess30m.toISOString(),
                  "format": "2006-01-02T15:04:05Z07:00"
                }
              }
            },
            !term
              ? { "match_all": {} }
              : {
                "query_string": {
                  "query": term
                }
              }
          ]
        }
      },
      "sort": [
        "-@timestamp"
      ],
      "from": page
        ? _from
        : 0,
      "size": page
        ? page.size || 20
        : 20,
      "aggs": {
        "histogram": {
          "date_histogram": {
            "field": "@timestamp",
            "calendar_interval": "",
            "fixed_interval": "30s"
          }
        }
      }
    }
  )
}

/**
 * 
 * @param {string} term 
 * @param {Object} page
 * default long page 20
 * {
 *  from: 0,
 *  size: 20,
 *  count:237
 * } 
 * quantity of pages: 12
 * @returns 
 */
export const searchEmail = async (term, page) => {
  try {
    const res = await axios.post(
      `${HOST_PATH}/es/mails/_search`,
      getBody(term, page),
      getAuth()
    )

    return {
      _emails: res.data.hits.hits,
      _total: res.data.hits.total?.value
    }  
  } catch (error) {
    throw new Error(error)
  }
  
}
