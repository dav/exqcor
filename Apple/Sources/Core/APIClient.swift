import Foundation

public struct ScriptSummary: Codable, Identifiable {
    public let id: Int
    public let title: String
    public let description: String?
}

public struct SectionSummary: Codable, Identifiable {
    public let id: Int
    public let name: String
    public let description: String?
    public let scriptId: Int

    enum CodingKeys: String, CodingKey {
        case id
        case name
        case description
        case scriptId = "script_id"
    }
}

public enum APIError: Error {
    case invalidURL
    case invalidResponse
}

public final class APIClient: ObservableObject {
    @Published public var scripts: [ScriptSummary] = []
    @Published public var sections: [SectionSummary] = []
    @Published public var errorMessage: String?

    public var baseURL: String
    public var token: String?

    public init(baseURL: String) {
        self.baseURL = baseURL
    }

    public func loadScripts() {
        request(path: "/api/scripts") { (result: Result<[ScriptSummary], Error>) in
            self.handle(result: result) { self.scripts = $0 }
        }
    }

    public func loadSections(scriptId: Int) {
        request(path: "/api/scripts/\(scriptId)/sections") { (result: Result<[SectionSummary], Error>) in
            self.handle(result: result) { self.sections = $0 }
        }
    }

    private func request<T: Decodable>(path: String, completion: @escaping (Result<T, Error>) -> Void) {
        guard let url = URL(string: baseURL)?.appendingPathComponent(path) else {
            completion(.failure(APIError.invalidURL))
            return
        }

        var request = URLRequest(url: url)
        request.addValue("application/json", forHTTPHeaderField: "Accept")
        if let token {
            request.addValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error {
                DispatchQueue.main.async { completion(.failure(error)) }
                return
            }

            guard let response = response as? HTTPURLResponse,
                  (200...299).contains(response.statusCode),
                  let data = data else {
                DispatchQueue.main.async { completion(.failure(APIError.invalidResponse)) }
                return
            }

            do {
                let decoded = try JSONDecoder().decode(T.self, from: data)
                DispatchQueue.main.async { completion(.success(decoded)) }
            } catch {
                DispatchQueue.main.async { completion(.failure(error)) }
            }
        }.resume()
    }

    private func handle<T>(result: Result<T, Error>, onSuccess: (T) -> Void) {
        switch result {
        case .success(let value):
            errorMessage = nil
            onSuccess(value)
        case .failure(let error):
            errorMessage = error.localizedDescription
        }
    }
}
