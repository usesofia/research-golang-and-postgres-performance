import http from 'k6/http';
import { sleep, check } from 'k6';
import exec from 'k6/execution';
import { createFinancialRecords } from './populate.js';
export const options = {
  vus: 100,
  duration: '15s',
  // duration: '60s',
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  const orgId = Math.max(1, (exec.vu.idInTest % 10) + 1);

  const response = http.get(`${BASE_URL}/organizations/${orgId}/financial-records/reports/cash-flow`);

  if (response.status !== 200) {
    console.log(`Failed to get cash flow for organization ${orgId}: ${response.status} ${response.body}`);
  }

  createFinancialRecords(orgId, 10);
}
