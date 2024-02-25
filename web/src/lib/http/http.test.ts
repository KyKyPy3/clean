import type { HttpRequest } from '@/src/main/http';
import { faker } from '@faker-js/faker';
import { vi } from 'vitest'
import { AxiosHttpClient } from './http';

export const mockHttpResponse = () => ({
  data: faker.string.sample(),
  status: faker.number.int(),
})

export const mockHttpRequest = (): HttpRequest => ({
  url: faker.internet.url(),
  method: faker.helpers.arrayElement(['get', 'post', 'put', 'delete']),
  body: faker.string.sample(),
  headers: faker.helpers.arrayElements(["headre1", "header2"])
})

const mocks = vi.hoisted(() => ({
  request: vi.fn(),
}));

vi.mock('axios', async (original: any) => {
  const actual: any = await vi.importActual("axios");

  return {
    default: {
      ...actual.default,
      request: mocks.request,
    }
  }
});

describe('AxiosHttpClient', () => {
  it('should call axios with correct values', async () => {
    const request = mockHttpRequest();
    let sut = new AxiosHttpClient();
    mocks.request.mockResolvedValueOnce(mockHttpResponse());

    await sut.request(request);

    expect(mocks.request).toHaveBeenCalledWith({
      url: request.url,
      data: request.body,
      headers: request.headers,
      method: request.method
    })
  });

  it('should return correct response', async () => {
    const request = mockHttpRequest();
    let sut = new AxiosHttpClient();
    const mockResp = mockHttpResponse();
    mocks.request.mockResolvedValueOnce(mockResp);

    const httpResponse = await sut.request(request);

    expect(httpResponse).toEqual({
      statusCode: mockResp.status,
      body: mockResp.data,
    })
  });

  it('should return correct error', async () => {
    const request = mockHttpRequest();
    let sut = new AxiosHttpClient();

    mocks.request.mockRejectedValueOnce({
      response: {
        status: 500,
        data: "error",
      }
    });

    const httpResponse = await sut.request(request);

    expect(httpResponse).toEqual({
      statusCode: 500,
      body: "error",
    })
  });
});