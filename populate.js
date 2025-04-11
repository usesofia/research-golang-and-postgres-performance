import http from 'k6/http';
import { sleep, check } from 'k6';
import exec from 'k6/execution';

export const options = {
  vus: 100,
  duration: '15s',
  // duration: '60s',
};

// Base URL for the API
const BASE_URL = 'http://localhost:8080';

// Function to generate a random tag name
function generateRandomTagName() {
  const adjectives = ['Red', 'Blue', 'Green', 'Yellow', 'Purple', 'Orange', 'Black', 'White', 'Pink', 'Brown'];
  const nouns = ['Cat', 'Dog', 'Bird', 'Fish', 'Lion', 'Tiger', 'Bear', 'Wolf', 'Fox', 'Deer'];
  const randomNum = Math.floor(Math.random() * 1000);
  const randomAdjective = adjectives[Math.floor(Math.random() * adjectives.length)];
  const randomNoun = nouns[Math.floor(Math.random() * nouns.length)];
  return `${randomAdjective} ${randomNoun} ${randomNum}`;
}

// Function to generate a random date between two years ago and now
function generateRandomDate() {
  const now = new Date();
  const twoYearsAgo = new Date();
  twoYearsAgo.setFullYear(now.getFullYear() - 2);
  
  const randomTime = twoYearsAgo.getTime() + Math.random() * (now.getTime() - twoYearsAgo.getTime());
  return new Date(randomTime).toISOString()
}

// Function to get all tags for an organization
function getAllTags(orgId) {
  for(let attempt = 0; attempt < 32; attempt++) {
    try {
      const response = http.get(`${BASE_URL}/organizations/${orgId}/tags`);
      
      if (response.status !== 200) {
        console.log(`Error getting tags for organization ${orgId}: API returned status ${response.status}`);
        continue;
      }
      
      const responseBody = JSON.parse(response.body);
      if (!responseBody.data || !Array.isArray(responseBody.data)) {
        console.log(`Error getting tags for organization ${orgId}: Invalid response format`);
        continue;
      }
      
      return responseBody.data;
    } catch (error) {
      console.log(`Error getting tags for organization ${orgId}: ${error}`);
    }
  }
  
  // Return empty array instead of throwing error
  console.log(`Failed to get tags for organization ${orgId} after multiple attempts, returning empty array`);
  return [];
}

// Function to select random tags from an array of tags
function selectRandomTags(tags, count) {
  const availableTags = tags;
  const selectedTags = [];
  const numTags = Math.min(count, availableTags.length);
  
  for (let i = 0; i < numTags; i++) {
    const randomIndex = Math.floor(Math.random() * availableTags.length);
    selectedTags.push(availableTags[randomIndex]);
    availableTags.splice(randomIndex, 1);
  }
  
  return selectedTags;
}

// Function to create a financial records
export function createFinancialRecords(orgId, count) {
  const tags = getAllTags(orgId);
  let payloads = [];

  for(let i = 0; i < count; i++) {
    const direction = Math.random() > 0.5 ? 'IN' : 'OUT';
    const amount = Math.floor(Math.random() * 10000) + 1;
    const dueDate = generateRandomDate();
    const numTags = Math.floor(Math.random() * 4); // 0 to 3 tags

    const payload = {
      direction,
      amount,
      dueDate,
      tags: selectRandomTags(tags, numTags)
    };

    payloads.push(payload);
  }

  // console.log({payloads});
  
  const response = http.post(`${BASE_URL}/organizations/${orgId}/financial-records/bulk`, JSON.stringify(payloads), {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (response.status !== 201) {
    console.log(`Failed to create financial records: ${response.status} ${response.body}`);
  }
}

// Function to create tags for an organization
function createTagsForOrganization(orgId) {
  const tagIds = [];
  const nTagsPerOrganization = 32;

  // Validate organization ID
  if (!orgId || orgId <= 0) {
    console.log(`Invalid organization ID: ${orgId}, skipping tag creation`);
    return tagIds;
  }
  
  for (let i = 0; i < nTagsPerOrganization; i++) {
    try {
      const numberOfTags = getNumberOfTags(orgId);

      if(numberOfTags >= nTagsPerOrganization) {
        break;
      }

      const tagName = generateRandomTagName();
      
      // Ensure tag name is not empty
      if (!tagName) {
        console.log(`Generated empty tag name, skipping`);
        continue;
      }

      const response = http.post(
        `${BASE_URL}/organizations/${orgId}/tags`,
        JSON.stringify({ name: tagName }),
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.status !== 201) {
        console.log(`Failed to create tag for organization ${orgId}: ${response.status} ${response.body}`);
        continue;
      }

      check(response, {
        'is status 201': (r) => r.status === 201,
      });
      
      if (response.status === 201) {
        try {
          const tagData = JSON.parse(response.body);
          if (tagData && tagData.id) {
            tagIds.push(tagData.id);
          }
        } catch (parseError) {
          console.log(`Error parsing tag response: ${parseError}`);
        }
      }
    } catch (error) {
      console.log(`Error in tag creation loop: ${error}`);
    }
  }
  
  return tagIds;
}

function getNumberOfTags(orgId) {
  for(let attempt = 0; attempt < 32; attempt++) {
    try {
      const response = http.get(`${BASE_URL}/organizations/${orgId}/tags`);
      
      if (response.status !== 200) {
        console.log(`Error getting number of tags for organization ${orgId}: API returned status ${response.status}`);
        continue;
      }
      
      const responseBody = JSON.parse(response.body);
      if (!responseBody.pagination || typeof responseBody.pagination.total_items !== 'number') {
        console.log(`Error getting number of tags for organization ${orgId}: Invalid response format`);
        continue;
      }
      
      return responseBody.pagination.total_items;
    } catch (error) {
      console.log(`Error getting number of tags for organization ${orgId}: ${error}`);
    }
  }
  
  // Return default value instead of throwing error
  console.log(`Failed to get number of tags for organization ${orgId} after multiple attempts, returning 0`);
  return 0;
}

export default function () {
  // Randomly select an organization ID (1-10)
  const orgId = Math.max(1, (exec.vu.idInTest % 10) + 1);

  createTagsForOrganization(orgId);
  createFinancialRecords(orgId, 400);
}
