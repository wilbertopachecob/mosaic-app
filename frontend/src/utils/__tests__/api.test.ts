import { parseJson } from "../api";

function mockResponse(body: string, init?: ResponseInit): Response {
  return new Response(body, init);
}

describe("parseJson", () => {
  it("parses valid JSON", async () => {
    const data = await parseJson<{ ok: boolean }>(
      mockResponse('{"ok":true}')
    );
    expect(data).toEqual({ ok: true });
  });

  it("returns null for an empty body", async () => {
    expect(await parseJson(mockResponse(""))).toBeNull();
  });

  it("returns null for invalid JSON", async () => {
    expect(await parseJson(mockResponse("not json"))).toBeNull();
  });
});
