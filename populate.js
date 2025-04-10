import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 100,
  duration: '15s',
  // duration: '60s',
};

// Base URL for the API
const BASE_URL = 'http://localhost:8080';

// Store tag IDs by organization
const organizationTags = {};

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

// Function to get random tags for an organization
function getRandomTags(orgId, count) {
  if (!organizationTags[orgId] || organizationTags[orgId].length === 0) {
    return [];
  }
  
  const availableTags = [...organizationTags[orgId]];
  const selectedTags = [];
  const numTags = Math.min(count, availableTags.length);
  
  for (let i = 0; i < numTags; i++) {
    const randomIndex = Math.floor(Math.random() * availableTags.length);
    selectedTags.push(availableTags[randomIndex]);
    availableTags.splice(randomIndex, 1);
  }
  
  return selectedTags;
}

// Function to create a financial record
function createFinancialRecord(orgId) {
  const direction = Math.random() > 0.5 ? 'IN' : 'OUT';
  const amount = Math.floor(Math.random() * 10000) + 1;
  const dueDate = generateRandomDate();
  const numTags = Math.floor(Math.random() * 4); // 0 to 3 tags
  const tags = getRandomTags(orgId, numTags);
  
  const response = http.post(
    `${BASE_URL}/organizations/${orgId}/financial-records`,
    JSON.stringify({
      direction,
      amount,
      dueDate,
      tags: tags.map(tagId => ({ id: tagId }))
    }),
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );

  check(response, {
    'is status 201': (r) => r.status === 201,
  });
  
  return response;
}

// Function to create tags for an organization
function createTagsForOrganization(orgId, count) {
  const tagIds = [];
  
  for (let i = 0; i < count; i++) {
    const tagName = generateRandomTagName();
    
    const response = http.post(
      `${BASE_URL}/organizations/${orgId}/tags`,
      JSON.stringify({ name: tagName }),
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );

    check(response, {
      'is status 201': (r) => r.status === 201,
    });
    
    if (response.status === 201) {
      const tagData = JSON.parse(response.body);
      tagIds.push(tagData.id);
    }
  }
  
  return tagIds;
}

export default function () {
  // Randomly select an organization ID (1-10)
  const orgId = Math.floor(Math.random() * 10) + 1;
  
  // Initialize organization tags if not already done
  if (!organizationTags[orgId]) {
    // Create 5 tags for this organization
    organizationTags[orgId] = createTagsForOrganization(orgId, 20);
  }
  
  // Create 32 financial records in parallel (chunks of 8)
  const chunks = 4; // 4 chunks of 8 records each
  const recordsPerChunk = 8;
  
  for (let chunk = 0; chunk < chunks; chunk++) {
    const promises = [];
    
    for (let i = 0; i < recordsPerChunk; i++) {
      promises.push(createFinancialRecord(orgId));
    }
    
    // Wait for all records in this chunk to be created
    Promise.all(promises);
    
    // Sleep for a short time between chunks
    sleep(0.5);
  }
  
  // Sleep for 1 second as requested
  sleep(1);
}
