import { JsonObject, JsonProperty } from 'typescript-json-serializer';

@JsonObject()
export class ScrapePropertiesJob {
  @JsonProperty() ID: number = 0
  @JsonProperty() CreatedAt: Date = new Date()
  @JsonProperty() UpdatedAt: Date = new Date()
  @JsonProperty() DeletedAt: Date | undefined = new Date()
  @JsonProperty() Url: string = ""
  @JsonProperty() Type: string = ""
  @JsonProperty() State: string = ""
  @JsonProperty() Progress: number = 0
  @JsonProperty() Message: string | undefined = ""
  @JsonProperty() Tag: string | undefined = ""
}
