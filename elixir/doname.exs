[name | _] = System.argv()

Application.ensure_all_started(:inets)
Application.ensure_all_started(:ssl)
{:ok, {_, _, body}} = :httpc.request("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
body
|> to_string()
|> String.trim()
|> String.split("\n")
|> Stream.drop(1)
|> Stream.map(fn tld -> String.replace(tld, ~r/--?.*/, "") end)
|> Stream.dedup()
|> Stream.map(&String.downcase/1)
|> Stream.filter(fn tld -> String.contains?(name, tld) and not String.starts_with?(name, tld) end)
|> Stream.map(fn tld ->
  [a, b] = String.split(name, tld, parts: 2)
  case b do
    "" -> "#{a}.#{tld}"
    b -> "#{a}.#{tld}/#{b}"
  end
end)
|> Enum.join("\n")
|> IO.puts()
