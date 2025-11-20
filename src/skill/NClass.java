package skill;

/*
 * @Author Coder: Nguyễn Tấn Tài
 * @Description: Ngọc Rồng Kiwi - Máy Chủ Chuẩn Teamobi 2025
 * @Group Zalo: https://zalo.me/g/toiyeuvietnam2025
 */
import player.system.Template.SkillTemplate;
import java.util.ArrayList;
import java.util.List;
import utils.Util;

public class NClass {

    public int classId;

    public String name;

    public List<SkillTemplate> skillTemplatess = new ArrayList<>();

    public SkillTemplate getSkillTemplate(int tempId) {
        for (SkillTemplate skillTemplate : skillTemplatess) {
            if (skillTemplate.id == tempId) {
                return skillTemplate;
            }
        }
        return null;
    }
public List<Skill> getSkills(int tempId) {
        for (SkillTemplate skillTemplate : skillTemplatess) {
            if (skillTemplate.id == tempId) {
                return skillTemplate.skillss;
            }
        }
        return null;
    }
    public SkillTemplate getSkillTemplateByName(String name) {
        for (SkillTemplate skillTemplate : skillTemplatess) {
            if ((Util.removeAccent(skillTemplate.name).toUpperCase()).contains((Util.removeAccent(name)).toUpperCase())) {
                return skillTemplate;
            }
        }
        return null;
    }

}
