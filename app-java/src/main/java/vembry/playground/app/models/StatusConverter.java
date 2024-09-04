
package vembry.playground.app.models;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter(autoApply = true)
public class StatusConverter implements AttributeConverter<Status, String> {

    @Override
    public String convertToDatabaseColumn(Status status) {
        if (status == null) {
            return null;
        }
        return status.name();
    }

    @Override
    public Status convertToEntityAttribute(String dbData) {
        if (dbData == null) {
            return null;
        }
        return Status.valueOf(dbData);
    }
}
