import request from '@/utils/request'

export default {
  /* curl --header "Content-Type: application/json" --request POST --data '{"params":["a","b"]}' http://host.docker.internal:4396/query1
  */
  query1(queryVO) {
    return request({
      url: `/query1`,
      method: 'post',
      data: queryVO
    })
  },
  query2(queryVO) {
    return request({
      url: `/query2`,
      method: 'post',
      data: queryVO
    })
  },
  query3(queryVO) {
    return request({
      url: `/query3`,
      method: 'post',
      data: queryVO
    })
  }

}
