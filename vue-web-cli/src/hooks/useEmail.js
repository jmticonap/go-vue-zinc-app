import { ref } from 'vue'
import axios from 'axios'
import { HOST_PATH } from '../../parameters'
import { getAuth } from '../../utils'

const getBody = term => {
  const currentDate = new Date()
  const dateLess30m = new Date(currentDate - (30 * 60 * 1000))
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
      "from": 0,
      "size": 100,
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

export const useEmail = async () => {
  const term = ref('')
  const data = ref({})
  const fetchData = async () => {
    data = await axios.post(
      `${HOST_PATH}/es/mails/_search`,
      getBody(term),
      getAuth()
    ).data
  }
  
  return {
    term,
    data,
    fetchData
  }
}
